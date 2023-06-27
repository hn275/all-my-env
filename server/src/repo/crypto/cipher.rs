#![allow(unused)]
use aes_gcm::{
    self,
    aead::{
        self,
        generic_array::{typenum::U32, GenericArray},
        Aead, AeadCore, KeyInit,
    },
    Aes256Gcm,
};
use hkdf::{self, Hkdf};
use sha2::Sha256;

pub type Nonce = [u8; 12];
pub type KeySecret = Vec<u8>;
pub type CipherResult = (Vec<u8>, Nonce);

#[derive(Debug)]
pub enum Error {
    InvalidKey(String),
    EncryptionFailed(String),
    InvalidNonce,
}

pub enum KeyType {
    RowKey,
}
impl KeyType {
    fn get_secret(&self) -> String {
        match self {
            Self::RowKey => self.get_env("ROW_KEY_SECRET"),
        }
    }

    fn get_env(&self, key: &str) -> String {
        use std::env::var;
        var(key)
            .expect(format!("`{}` not set", key).as_str())
            .to_owned()
    }
}

pub trait KeyGen {
    fn key(&self) -> Vec<u8>;
}

pub struct Key(KeySecret);
impl Key {
    pub fn new(k: KeyType) -> Key {
        let key = match k {
            KeyType::RowKey => k.get_secret(),
        };
        Key(key.as_bytes().to_owned())
    }

    pub fn generate_key(&self, c: &impl KeyGen) -> Result<Cipher, Error> {
        use aes_gcm::Key;

        let salt = None;
        let master_key = c.key();
        let kdf = Hkdf::<Sha256>::new(salt, &master_key);

        let mut output_key = [0u8; 32];
        kdf.expand(&self.0, &mut output_key)
            .map_err(|err| Error::InvalidKey(err.to_string()))?;

        let key = Key::<Aes256Gcm>::from_slice(output_key.as_ref()).to_owned();

        return Ok(Cipher(key));
    }
}

pub struct Cipher(GenericArray<u8, U32>);
impl Cipher {
    pub fn seal(&self, plaintext: &[u8]) -> Result<CipherResult, Error> {
        let nonce = Aes256Gcm::generate_nonce(&mut aead::OsRng);

        let cipher = aes_gcm::Aes256Gcm::new(&self.0);
        let ciphertext = cipher
            .encrypt(&nonce, plaintext)
            .map_err(|err| Error::EncryptionFailed(err.to_string()))?;

        let nonce: Nonce = nonce.try_into().map_err(|_| Error::InvalidNonce)?;

        return Ok((ciphertext, nonce));
    }

    pub fn open(&self, ciphertext: &[u8], n: Nonce) -> Result<Vec<u8>, Error> {
        let nonce = GenericArray::from_slice(&n);
        let cipher = Aes256Gcm::new(&self.0);

        let plaintext = Aes256Gcm::new(&self.0)
            .decrypt(&nonce, ciphertext)
            .map_err(|err| Error::EncryptionFailed(err.to_string()))?;

        return Ok(plaintext);
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::env::set_var;

    #[derive(Clone, Debug)]
    struct Foo {
        id: String,
        key: String,
        value: String,
    }

    impl KeyGen for Foo {
        fn key(&self) -> Vec<u8> {
            [self.id.as_bytes(), self.key.as_bytes()].concat()
        }
    }

    #[test]
    fn cipher_decipher() {
        set_var("ROW_KEY_SECRET", "testing");

        let foo = Foo {
            id: "123".to_owned(),
            key: "foo".to_owned(),
            value: "bar".to_owned(),
        };

        let (ciphertext, nonce) = Key::new(KeyType::RowKey)
            .generate_key(&foo)
            .unwrap()
            .seal(foo.value.as_bytes())
            .unwrap();
        assert_ne!(foo.value.as_bytes(), ciphertext);

        let opened = Key::new(KeyType::RowKey)
            .generate_key(&foo)
            .unwrap()
            .open(ciphertext.as_slice(), nonce)
            .unwrap();
        assert_eq!(foo.value.as_bytes(), opened);

        let value = std::str::from_utf8(opened.as_slice()).unwrap();
        assert_eq!(value, "bar");
    }
}

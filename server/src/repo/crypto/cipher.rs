use aes_gcm::{
    self,
    aead::{self, generic_array::GenericArray, Aead, AeadCore, KeyInit},
};
use hkdf::{self, Hkdf};
use sha2::Sha256;

#[derive(Debug)]
pub enum Error {
    ErrInvalidKey(String),
    ErrEncryption(String),
    ErrInvalidNonce,
}

pub enum KeyType {
    RowKey,
}

impl KeyType {
    fn get_secret(&self) -> String {
        match self {
            Self::RowKey => get_env("ROW_KEY_SECRET"),
        }
    }
}

fn get_env(key: &str) -> String {
    use std::env::var;
    var(key)
        .expect(format!("`{}` not set", key).as_str())
        .to_owned()
}

pub trait CipherComponent {
    fn seal(&mut self, c: CipherData);
    fn open(&self) -> CipherData;
    fn master_key(&self) -> Vec<u8>;
}

#[derive(Debug, serde::Serialize)]
pub struct CipherData {
    pub ciphertext: Vec<u8>,
    pub nonce: Option<Nonce>, // nonce size is 12 from aesgcm
}

pub type Nonce = [u8; 12];
pub struct Cipher {
    secret: Vec<u8>,
}
impl Cipher {
    pub fn new(key_type: KeyType) -> Cipher {
        Cipher {
            secret: key_type.get_secret().as_bytes().to_vec(),
        }
    }

    fn generate_key(&self, c: &impl CipherComponent) -> Result<[u8; 32], Error> {
        let salt = None;
        let kdf = Hkdf::<Sha256>::new(salt, c.master_key().as_ref());

        let mut output_key = [0u8; 32];
        kdf.expand(&self.secret, &mut output_key)
            .map_err(|err| Error::ErrInvalidKey(err.to_string()))?;

        return Ok(output_key);
    }

    pub fn seal(&mut self, c: &mut impl CipherComponent) -> Result<(), Error> {
        let key = self.generate_key(c)?;
        let key = aes_gcm::Key::<aes_gcm::Aes256Gcm>::from_slice(key.as_ref()).to_owned();

        let nonce = match &c.open().nonce {
            None => aes_gcm::Aes256Gcm::generate_nonce(&mut aead::OsRng),
            Some(n) => GenericArray::from_slice(n).to_owned(),
        };

        let cipher = aes_gcm::Aes256Gcm::new(&key);

        let ciphertext = cipher
            .encrypt(&nonce, c.master_key().as_ref())
            .map_err(|err| Error::ErrEncryption(err.to_string()))?;

        let nonce: [u8; 12] = nonce.try_into().map_err(|_| Error::ErrInvalidNonce)?;
        let cipher_data = CipherData {
            ciphertext,
            nonce: Some(nonce),
        };
        c.seal(cipher_data);

        return Ok(());
    }

    pub fn open(&self, c: &impl CipherComponent) -> Result<Vec<u8>, Error> {
        let cipher_data = c.open();
        let nonce = match &cipher_data.nonce {
            None => Err(Error::ErrInvalidNonce)?,
            Some(n) => GenericArray::from_slice(n).to_owned(),
        };

        let key = self.generate_key(c)?;
        let key = aes_gcm::Key::<aes_gcm::Aes256Gcm>::from_slice(key.as_ref()).to_owned();

        let cipher = aes_gcm::Aes256Gcm::new(&key);

        let plaintext = cipher
            .decrypt(&nonce, cipher_data.ciphertext.as_slice())
            .map_err(|err| Error::ErrEncryption(err.to_string()))?;

        return Ok(plaintext);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[derive(Clone, Debug)]
    struct Foo {
        id: String,
        key: String,
        value: String,
        nonce: Option<[u8; 12]>,
    }

    impl CipherComponent for Foo {
        fn open(&self) -> CipherData {
            CipherData {
                ciphertext: self.value.as_bytes().to_vec(),
                nonce: self.nonce,
            }
        }

        fn seal(&mut self, c: CipherData) {
            self.value = hex::encode(c.ciphertext);
            self.nonce = c.nonce
        }

        fn master_key(&self) -> Vec<u8> {
            self.id.as_bytes().to_owned()
        }
    }

    #[test]
    fn it_works() {
        let mut foo = Foo {
            id: "123".to_owned(),
            key: "foo".to_owned(),
            value: "bar".to_owned(),
            nonce: None,
        };

        let foo2 = foo.clone();
        Cipher::new(KeyType::RowKey).seal(&mut foo).unwrap();

        assert_ne!(foo.value, foo2.value);
        assert!(foo.nonce.is_some());
        assert_eq!(foo.id, foo2.id);
        assert_eq!(foo.key, foo2.key);

        Cipher::new(KeyType::RowKey).open(&mut foo).unwrap();

        // assert_eq!(foo.value, foo2.value);
        // assert!(foo.nonce.is_some());
        // assert_eq!(foo.id, foo2.id);
        // assert_eq!(foo.key, foo2.key);
    }
}

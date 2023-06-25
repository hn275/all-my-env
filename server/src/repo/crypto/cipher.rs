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
    fn input_key(&self) -> &[u8];
    fn input_data(&self) -> &[u8];

    // returns data used to open
    fn open(&self) -> CipherData;
    // returns nonce (if any)
    fn nonce(&self) -> Option<Vec<u8>>;
    // force updating the struct once Cipher data is generated
    fn update(&mut self, d: CipherData);
}

#[derive(Debug, serde::Serialize)]
pub struct CipherData {
    pub ciphertext: Vec<u8>,
    pub nonce: [u8; 12], // nonce size is 12 from aesgcm
}

pub struct Cipher(String);
impl Cipher {
    pub fn new(key_type: KeyType) -> Cipher {
        Cipher(key_type.get_secret())
    }

    fn generate_key(&self, buf: &impl CipherComponent) -> Result<[u8; 32], Error> {
        let salt = None;
        let key = buf.input_key();
        let kdf = Hkdf::<Sha256>::new(salt, key);

        let secret = self.0.as_bytes();
        let mut output_key = [0u8; 32];
        kdf.expand(&secret, &mut output_key)
            .map_err(|err| Error::ErrInvalidKey(err.to_string()))?;

        return Ok(output_key);
    }

    pub fn seal(&self, buf: &mut impl CipherComponent) -> Result<(), Error> {
        let key = self.generate_key(buf)?;
        let key = aes_gcm::Key::<aes_gcm::Aes256Gcm>::from_slice(&key);

        let nonce = match &buf.nonce() {
            None => aes_gcm::Aes256Gcm::generate_nonce(&mut aead::OsRng),
            Some(n) => {
                let n: [u8; 12] = n
                    .to_owned()
                    .try_into()
                    .map_err(|_| Error::ErrInvalidNonce)?;

                GenericArray::from_slice(&n).to_owned()
            }
        };

        let cipher = aes_gcm::Aes256Gcm::new(&key);
        let plaintext = buf.input_data();

        let ciphertext = cipher
            .encrypt(&nonce, plaintext)
            .map_err(|err| Error::ErrEncryption(err.to_string()))?;

        let nonce: [u8; 12] = nonce.try_into().map_err(|_| Error::ErrInvalidNonce)?;
        let cipher_data = CipherData { ciphertext, nonce };
        buf.update(cipher_data);

        return Ok(());
    }

    pub fn open(&self, buf: &impl CipherComponent) -> Result<Vec<u8>, Error> {
        let key = self.generate_key(buf)?;
        let key = aes_gcm::Key::<aes_gcm::Aes256Gcm>::from_slice(&key);
        let cipher_data = buf.open();

        let cipher = aes_gcm::Aes256Gcm::new(&key);
        let nonce = GenericArray::from(cipher_data.nonce);
        let plaintext = cipher
            .decrypt(&nonce, cipher_data.ciphertext.as_ref())
            .map_err(|err| Error::ErrEncryption(err.to_string()))?;

        return Ok(plaintext);
    }
}

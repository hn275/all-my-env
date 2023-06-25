use aes_gcm::aead::{self, generic_array::GenericArray, Aead, AeadCore, KeyInit};
use hkdf::{self, Hkdf, InvalidLength};
use sha2::Sha256;

pub enum KeyType {
    RowKey,
}

impl KeyType {
    fn get_key(&self) -> String {
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
    fn key(&self) -> Vec<u8>;
    fn data(&self) -> &[u8];
    fn seal(&mut self, c: CipherData);
    fn open(&self) -> CipherData;
}

pub struct CipherKey([u8; 32]);
impl CipherKey {
    pub fn new(key: KeyType, buf: &impl CipherComponent) -> Result<CipherKey, InvalidLength> {
        let salt = None;
        let kdf = Hkdf::<Sha256>::new(salt, &buf.key());

        let secret = key.get_key().as_bytes().to_owned();
        let mut output_key = [0u8; 32];
        kdf.expand(&secret, &mut output_key)?;

        return Ok(CipherKey(output_key));
    }

    pub fn seal(&self, buf: &mut impl CipherComponent) -> Result<(), aes_gcm::Error> {
        let key: &aes_gcm::Key<aes_gcm::Aes256Gcm> = &self.0.into();
        let cipher = aes_gcm::Aes256Gcm::new(&key);
        let nonce = aes_gcm::Aes256Gcm::generate_nonce(&mut aead::OsRng);

        let ciphertext = cipher.encrypt(&nonce, buf.data())?;
        let nonce: [u8; 12] = nonce.try_into().expect("handle this error");

        let cipher_data = CipherData { ciphertext, nonce };
        buf.seal(cipher_data);

        return Ok(());
    }

    pub fn open(&self, buf: &impl CipherComponent) -> Result<Vec<u8>, ()> {
        let cipher_data = buf.open();
        let key = aes_gcm::Key::<aes_gcm::Aes256Gcm>::from_slice(&self.0);
        let cipher = aes_gcm::Aes256Gcm::new(&key);
        let nonce = GenericArray::from(cipher_data.nonce);
        let plaintext = cipher
            .decrypt(&nonce, cipher_data.ciphertext.as_ref())
            .unwrap();

        return Ok(plaintext);
    }
}

#[derive(Debug, serde::Serialize)]
pub struct CipherData {
    pub ciphertext: Vec<u8>,
    pub nonce: [u8; 12], // nonce size is 12 from aesgcm
}

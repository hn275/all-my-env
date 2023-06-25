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

pub trait KeyComponent {
    fn to_bytes(&self) -> Vec<u8>;
}

pub fn generate(key: KeyType, buf: &impl KeyComponent) -> Result<[u8; 32], InvalidLength> {
    let salt = None;
    let kdf = Hkdf::<Sha256>::new(salt, &buf.to_bytes());

    let secret = key.get_key().as_bytes().to_owned();
    let mut output_key = [0u8; 32];
    kdf.expand(&secret, &mut output_key)?;

    return Ok(output_key);
}

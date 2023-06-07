use aes_gcm::{Aes256Gcm, Key};

pub enum CipherKeys {
    Auth,
}

pub struct Cipher<'a> {
    value: &'a str,
    key: Option<&'a str>,
}

impl<'a> Cipher<'a> {
    pub fn new(value: &str) -> Cipher {
        return Cipher { value, key: None };
    }

    pub fn key(&mut self, key: CipherKeys) -> &'a mut Cipher {
        let k = match key {
            CipherKeys::Auth => "hSchE34mgOHmvsQokZGSZY4jPQbqD9qY",
        };
        self.key = Some(k);
        return self;
    }

    pub fn decrypt() -> Result<(), ()> {
        return Ok(());
    }
}

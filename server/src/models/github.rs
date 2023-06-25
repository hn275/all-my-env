use crate::repo::crypto::cipher;
use std::collections::HashMap;

pub struct Repo {
    pub write_access: AccessControl,
    pub meta: Meta,
    pub variables: Vec<Variable>,
}

pub type UserID = String;
pub type AccessControl = HashMap<UserID, User>;
pub struct User {
    pub username: String,
    pub name: String,
}

pub struct Meta {
    pub id: String,
    pub repo_name: String,
    pub repo_id: String,
    pub owner: String,
}

pub struct Variable {
    pub key: String,
    pub value: String, // base64 encoding
    pub nonce: String, // base64 encoding
}

impl cipher::CipherComponent for Variable {
    fn open(&self) -> cipher::CipherData {
        todo!();
    }

    fn nonce(&self) -> Option<Vec<u8>> {
        todo!();
    }

    fn update(&mut self, _d: cipher::CipherData) {
        todo!();
    }

    fn input_key(&self) -> &[u8] {
        todo!();
    }

    fn input_data(&self) -> &[u8] {
        todo!();
    }
}

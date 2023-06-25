use std::collections::HashMap;

pub struct Repo {
    pub write_access: AccessControl,
    pub meta: Meta,
    pub variables: Vec<Variable>,
}

pub type user_id = String;
pub type AccessControl = HashMap<user_id, User>;
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
    pub value: String,
    pub nonce: String, // hex encoding of nonce
}

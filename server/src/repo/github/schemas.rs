use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GithubAccount {
    #[serde(skip_deserializing)]
    pub token: String,
    pub login: String,
    pub id: u32,
    #[serde(rename(serialize = "avatarUrl"))]
    pub avatar_url: String,
    pub name: String,
    pub email: String,
    pub bio: String,
    #[serde(rename(serialize = "publicRepos"))]
    pub public_repos: u32,
    pub followers: u32,
    pub following: u32,
    #[serde(rename(serialize = "createdAt"))]
    pub created_at: String,
    #[serde(rename(serialize = "updatedAt"))]
    pub updated_at: String,
}

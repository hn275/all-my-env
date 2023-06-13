use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GithubAccount {
    pub login: String,
    pub id: i64,
    #[serde(rename(serialize = "avatarUrl"))]
    pub avatar_url: String,
    pub name: String,
    pub email: String,
    pub bio: String,
    #[serde(rename(serialize = "publicRepos"))]
    pub public_repos: i64,
    pub followers: i64,
    pub following: i64,
    #[serde(rename(serialize = "createdAt"))]
    pub created_at: String,
    #[serde(rename(serialize = "updatedAt"))]
    pub updated_at: String,
}

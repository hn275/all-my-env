use crate::repo::error::ApiError;
use actix_web::http;
use reqwest;
use serde::{Deserialize, Serialize};
use std::env;

pub struct AuthClient<'a> {
    pub client_id: String,
    pub client_secret: String,
    pub code: &'a str,
    pub redirect_uri: Option<&'a str>,
    pub scope: Vec<&'a str>,
}

#[derive(Serialize, Deserialize)]
pub struct AuthToken {
    access_token: String,
    scope: String,
    token_type: String,
}

impl<'a> AuthClient<'a> {
    pub fn new(code: &'a str) -> AuthClient {
        let client_id = env::var("GITHUB_CLIENT_ID").expect("`GITHUB_CLIENT_ID` not set");
        let client_secret =
            env::var("GITHUB_CLIENT_SECRET").expect("`GITHUB_CLIENT_SECRET` not set");

        return AuthClient {
            client_id,
            client_secret,
            code,
            redirect_uri: None,
            scope: vec!["user:read repo:read"],
        };
    }

    pub async fn get_auth_token(&self) -> Result<String, ApiError> {
        let query_params = format!(
            "client_id={}&client_secret={}&code={}",
            self.client_id, self.client_secret, self.code
        );

        let url = format!(
            "https://github.com/login/oauth/access_token?{}",
            query_params
        );

        let response = reqwest::Client::new()
            .post(url)
            .header("Content-Type", "application/json")
            .send()
            .await
            .map_err(|err| {
                eprintln!("{}", err);
                return ApiError {
                    code: http::StatusCode::BAD_GATEWAY,
                    message: None,
                };
            })?;

        if let http::StatusCode::OK = response.status() {
            let result = response.text().await.map_err(|err| {
                eprintln!("{}", err);
                return ApiError {
                    code: http::StatusCode::INTERNAL_SERVER_ERROR,
                    message: None,
                };
            })?;

            let token = result.split('&').collect::<Vec<&str>>()[0];
            let result = token.split('=').collect::<Vec<&str>>();

            if result[0] == "error" {
                let err = ApiError {
                    code: http::StatusCode::BAD_REQUEST,
                    message: Some(result[1].to_owned()),
                };
                return Err(err);
            }

            return Ok(String::from(result[1]));
        } else {
            let err = ApiError {
                code: http::StatusCode::BAD_GATEWAY,
                message: None,
            };

            return Err(err);
        }
    }
}

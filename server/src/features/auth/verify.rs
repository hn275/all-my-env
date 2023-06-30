use crate::repo::{
    error::ApiError,
    github::{client::GithubClient, schemas::GithubAccount},
};
use actix_web::{
    http::{header, StatusCode},
    web, HttpResponse, HttpResponseBuilder,
};
use hmac::{Hmac, Mac};
use jwt::SignWithKey;
use reqwest;
use serde::{Deserialize, Serialize};
use sha2::Sha256;
use std::future::Future;
use std::{env::var, error::Error};

#[derive(Deserialize, Serialize, Debug)]
pub struct Token {
    pub code: String,
}

pub async fn verify_code(token: web::Json<Token>) -> Result<HttpResponse, ApiError> {
    let response = get_oauth(&token.into_inner().code).await.map_err(|err| {
        return ApiError::bad_gateway(err.to_string());
    })?;

    let response_status = response.status();
    if StatusCode::OK != response_status {
        let msg = response.text().await.map_err(|err| {
            return ApiError::internal(Box::new(err));
        })?;
        return Err(ApiError::new(response_status, Some(msg)));
    }

    let payload = response.text().await.map_err(|err| {
        return ApiError::internal(Box::new(err));
    })?;

    let bearer_token = get_token(payload.as_ref())?;
    let mut github_account = get_user_account(bearer_token).await?;
    github_account.token = String::from(bearer_token);

    let code = sign_jwt(&github_account).map_err(|err| ApiError::internal(err))?;
    let auth_token = Token { code };

    let response = HttpResponseBuilder::new(StatusCode::OK)
        .insert_header(header::ContentType::json())
        .json(&auth_token);

    return Ok(response);
}

fn sign_jwt(u: &GithubAccount) -> Result<String, Box<dyn Error>> {
    let secret = var("JWT_SECRET").expect("`JWT_SECRET` not set");
    let key: Hmac<Sha256> = Hmac::new_from_slice(&secret.as_bytes())?;
    let token = u.sign_with_key(&key)?;
    return Ok(token);
}

async fn get_user_account(token: &str) -> Result<GithubAccount, ApiError> {
    let client = GithubClient::new_with_token(token);

    let response = client.get("/user").send().await.map_err(|err| {
        return ApiError::internal(Box::new(err));
    })?;

    let status_code = response.status();

    if status_code != StatusCode::OK {
        let msg = response.text().await.map_err(|err| {
            return ApiError::internal(Box::new(err));
        })?;
        return Err(ApiError::new(status_code, Some(msg)));
    }

    let github_account = response.json::<GithubAccount>().await.map_err(|err| {
        return ApiError::internal(Box::new(err));
    })?;

    return Ok(github_account);
}

fn get_token(payload: &str) -> Result<&str, ApiError> {
    let result = payload.split('&').collect::<Vec<&str>>()[0];
    let payload = result.split('=').collect::<Vec<&str>>();
    if payload[0] == "error" {
        let message = Some(String::from(payload[1]));
        return Err(ApiError::unauthorized(message));
    }

    return Ok(payload[1]);
}

fn get_oauth(code: &str) -> impl Future<Output = Result<reqwest::Response, reqwest::Error>> {
    let client_id = var("GITHUB_CLIENT_ID").expect("`GITHUB_CLIENT_ID` not set");
    let client_secret = var("GITHUB_CLIENT_SECRET").expect("`GITHUB_CLIENT_SECRET` not set");

    let query_params = format!(
        "client_id={}&client_secret={}&code={}",
        client_id, client_secret, code
    );

    let url = format!(
        "https://github.com/login/oauth/access_token?{}",
        query_params
    );

    return reqwest::Client::new()
        .post(url)
        .header("Content-Type", "application/json")
        .send();
}

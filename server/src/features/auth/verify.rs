use std::future::Future;

use crate::repo::{error::ApiError, github::GithubClient};
use actix_web::{
    http::{header, StatusCode},
    web, HttpResponse, HttpResponseBuilder,
};
use reqwest;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct Token {
    pub code: String,
}

pub async fn verify_code(token: web::Json<Token>) -> Result<HttpResponse, ApiError> {
    let response = get_token(&token.into_inner().code).await.map_err(|err| {
        eprintln!("{}", err);
        return ApiError::new(StatusCode::BAD_GATEWAY, Some(err.to_string()));
    })?;

    if StatusCode::OK != response.status() {
        return Err(ApiError::new(response.status(), None));
    }

    let payload = response.text().await.map_err(|err| {
        eprintln!("{}", err);
        return ApiError::new(StatusCode::INTERNAL_SERVER_ERROR, None);
    })?;

    let result = payload.split('&').collect::<Vec<&str>>()[0];
    let payload = result.split('=').collect::<Vec<&str>>();

    if payload[0] == "error" {
        return Err(ApiError {
            code: StatusCode::BAD_REQUEST,
            message: Some(payload[1].to_owned()),
        });
    }

    let token = payload[1];

    let client = GithubClient::new_with_token(token);
    let result = client
        .get("/user")
        .send()
        .await
        .unwrap()
        .text()
        .await
        .unwrap();

    println!("{}", result);

    let response = HttpResponseBuilder::new(StatusCode::OK)
        .insert_header(header::ContentType::plaintext())
        // .body(result);
        .body(payload[1].to_owned());

    return Ok(response);
}

fn get_token(code: &str) -> impl Future<Output = Result<reqwest::Response, reqwest::Error>> {
    use std::env::var;

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

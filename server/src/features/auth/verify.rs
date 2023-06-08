use crate::repo::{error::ApiError, github};
use actix_web::{
    http::{header, StatusCode},
    web, HttpResponse, HttpResponseBuilder,
};
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct Token {
    pub code: String,
}

pub async fn verify_code(token: web::Json<Token>) -> Result<HttpResponse, ApiError> {
    let result = github::AuthClient::new(&token.into_inner().code)
        .get_auth_token()
        .await?;

    println!("{}", result);

    let response = HttpResponseBuilder::new(StatusCode::OK)
        .insert_header(header::ContentType::plaintext())
        .body(result);

    return Ok(response);
}

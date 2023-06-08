use actix_web;
use actix_web::{web, HttpResponse};
use reqwest;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct Token {
    pub code: String,
}

pub async fn verify_code(token: web::Json<Token>) -> impl actix_web::Responder {
    let t = get_access_token(&token.into_inner()).await.unwrap();
    dbg!("{}", &t);
    return HttpResponse::Ok();
}

async fn get_access_token(code: &Token) -> Result<String, ()> {
    let client = reqwest::Client::new();
    let response = client.post("https://github.com/login/oauth/access_token");
    let access_token = response
        .json(&code)
        .header("Content-Type", "application/json")
        .send()
        .await
        .map_err(|err| {
            eprintln!("request failed {}", err);
            return ();
        })?
        //.json::<Token>()
        .text()
        .await
        .map_err(|err| {
            eprintln!("struct marshal failed {}", err);
            return ();
        })?;

    Ok(access_token)
}

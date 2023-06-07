use actix_web;
use actix_web::{web, HttpResponse};
use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Token {
    pub code: String,
}

pub async fn verify_code(token: web::Json<Token>) -> impl actix_web::Responder {
    dbg!("{}", &token.code);
    return HttpResponse::Ok();
}

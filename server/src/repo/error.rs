use actix_web::{body, http, HttpResponse, HttpResponseBuilder, Responder};
use serde;

pub struct ApiError<'a> {
    code: http::StatusCode,
    message: Option<&'a str>,
}

#[derive(serde::Serialize)]
struct ErrorMessage<'a> {
    message: &'a str,
}

impl<'a> Responder for ApiError<'a> {
    type Body = body::BoxBody;
    fn respond_to(self, _req: &actix_web::HttpRequest) -> HttpResponse<Self::Body> {
        let mut responder = HttpResponseBuilder::new(self.code);
        match self.message {
            None => responder.finish(),
            Some(msg) => responder
                .insert_header(http::header::ContentType::json())
                .json(ErrorMessage { message: msg }),
        }
    }
}

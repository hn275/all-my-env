use actix_web::{body, http, HttpResponse, HttpResponseBuilder, ResponseError};
use serde;

#[derive(Debug)]
pub struct ApiError {
    pub code: http::StatusCode,
    pub message: Option<String>,
}

impl std::fmt::Display for ApiError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "")
    }
}

#[derive(serde::Serialize)]
struct ErrorMessage {
    error: String,
}

impl ResponseError for ApiError {
    fn status_code(&self) -> http::StatusCode {
        self.code
    }

    fn error_response(&self) -> HttpResponse<body::BoxBody> {
        let mut responder = HttpResponseBuilder::new(self.code);
        match &self.message {
            None => responder.finish(),
            Some(message) => responder
                .insert_header(http::header::ContentType::json())
                .json(ErrorMessage {
                    error: message.to_owned(),
                }),
        }
    }
}

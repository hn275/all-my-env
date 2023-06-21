use actix_web::{
    body,
    http::{self, StatusCode},
    HttpResponse, HttpResponseBuilder, ResponseError,
};
use serde;

#[derive(Debug)]
pub struct ApiError {
    pub code: StatusCode,
    pub message: Option<String>,
}

impl ApiError {
    pub fn new(code: StatusCode, message: Option<String>) -> ApiError {
        return ApiError { code, message };
    }

    #[allow(unused)] // TODO: remove this flag when on use
    pub fn bad_request(message: &str) -> ApiError {
        return ApiError {
            code: StatusCode::BAD_GATEWAY,
            message: Some(message.to_string()),
        };
    }

    #[allow(unused)] // TODO: remove this flag when on use
    pub fn not_found(message: Option<String>) -> ApiError {
        return ApiError {
            code: StatusCode::NOT_FOUND,
            message,
        };
    }

    #[allow(unused)] // TODO: remove this flag when on use
    pub fn forbidden() -> ApiError {
        return ApiError {
            code: StatusCode::FORBIDDEN,
            message: None,
        };
    }

    pub fn unauthorized(message: Option<String>) -> ApiError {
        return ApiError {
            code: StatusCode::UNAUTHORIZED,
            message,
        };
    }

    pub fn internal(err: Box<dyn std::error::Error>) -> ApiError {
        eprintln!("{}", err.to_string());
        return ApiError {
            code: StatusCode::INTERNAL_SERVER_ERROR,
            message: None,
        };
    }

    pub fn bad_gateway(message: String) -> ApiError {
        return ApiError {
            code: StatusCode::INTERNAL_SERVER_ERROR,
            message: Some(message),
        };
    }
}

impl std::fmt::Display for ApiError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "code: {}", &self.code)
    }
}

#[derive(serde::Serialize)]
struct ErrorMessage {
    error: String,
}

impl ResponseError for ApiError {
    fn status_code(&self) -> StatusCode {
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

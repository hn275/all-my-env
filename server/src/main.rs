use actix_cors::Cors;
use actix_web::{middleware::Logger, web, App, HttpServer};
use dotenv::dotenv;
use env_logger;

mod features;
mod repo;

use features::auth;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(env_logger::Env::default().default_filter_or("INFO"));

    let port = ("127.0.0.1", 8080);
    println!("serving at {}:{}", port.0, port.1);
    return HttpServer::new(|| {
        let cors = Cors::default().allow_any_origin().allow_any_header();

        App::new()
            .wrap(Logger::new("%a %s %b"))
            .wrap(cors)
            .route("/auth", web::post().to(auth::verify::verify_code))
            .service(test_route)
    })
    .bind(port)?
    .run()
    .await;
}

#[actix_web::get("/test")]
async fn test_route() -> impl actix_web::Responder {
    use repo::crypto::keygen;
    let components = vec!["foo".to_owned(), "bar".to_owned(), "baz".to_owned()];
    let key = keygen::generate(&components).unwrap();
    return hex::encode(key);
}

use actix_cors::Cors;
use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use env_logger;

mod features;
mod repo;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::default().default_filter_or("INFO"));

    let port = ("127.0.0.1", 8080);
    println!("serving at {}:{}", port.0, port.1);
    return HttpServer::new(|| {
        let cors = Cors::default().allow_any_origin().allow_any_header();

        App::new()
            .wrap(cors)
            .wrap(Logger::new("%a %s %b"))
            .route("/auth", web::post().to(features::auth::verify::verify_code))
    })
    .bind(port)?
    .run()
    .await;
}

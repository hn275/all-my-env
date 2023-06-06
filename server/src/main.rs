mod features;
mod repo;

use actix_web::{App, HttpServer};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    return HttpServer::new(|| {
        App::new()
        //
        //
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await;
}

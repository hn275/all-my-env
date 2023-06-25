use actix_cors::Cors;
use actix_web::{middleware::Logger, web, App, HttpServer};
use dotenv::dotenv;
use env_logger;

mod features;
mod models;
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

// example
#[derive(Debug)]
struct Foo {
    bar: String,
    cipher: String,
    nonce: Option<String>,
}

impl repo::crypto::keygen::CipherComponent for Foo {
    fn key(&self) -> Vec<u8> {
        self.bar.as_bytes().to_owned()
    }

    fn data(&self) -> &[u8] {
        self.cipher.as_bytes()
    }

    fn seal(&mut self, c: repo::crypto::keygen::CipherData) {
        self.cipher = hex::encode(c.ciphertext);
        self.nonce = Some(hex::encode(c.nonce));
    }

    fn open(&self) -> repo::crypto::keygen::CipherData {
        if let None = self.nonce {
            panic!("alskfjklsdjf");
        }

        let nonce = &self.nonce.as_ref().unwrap();
        repo::crypto::keygen::CipherData {
            nonce: hex::decode(nonce).unwrap().try_into().unwrap(),
            ciphertext: hex::decode(&self.cipher).unwrap(),
        }
    }
}

#[actix_web::get("/test")]
async fn test_route() -> impl actix_web::Responder {
    use repo::crypto::keygen;
    let mut foo = Foo {
        bar: "bar".to_owned(),
        cipher: "baz".to_owned(),
        nonce: None,
    };

    dbg!("before seal: {}", &foo);

    keygen::CipherKey::new(keygen::KeyType::RowKey, &foo)
        .unwrap()
        .seal(&mut foo)
        .unwrap();

    dbg!("after seal: {}", &foo);

    keygen::CipherKey::new(keygen::KeyType::RowKey, &foo)
        .unwrap()
        .open(&mut foo)
        .unwrap();

    dbg!("after open: {}", &foo);

    return std::str::from_utf8(&foo.bar.as_ref()).unwrap().to_owned();
}

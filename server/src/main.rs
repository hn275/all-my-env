use actix_cors::Cors;
use actix_web::{middleware::Logger, web, App, HttpServer};
use dotenv::dotenv;
use env_logger;

mod features;
// mod models;
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
        // .service(test_route)
    })
    .bind(port)?
    .run()
    .await;
}

// example
// NOTE: base64 encoding would be more size efficient,
// just doing hex encoding here bc going fast
// #[derive(Debug)]
// struct Foo {
//     bar: String,
//     cipher: String,
//     nonce: Option<String>,
// }
//
// impl repo::crypto::cipher::CipherComponent for Foo {
//     fn input_key(&self) -> &[u8] {
//         self.bar.as_bytes()
//     }
//
//     fn input_data(&self) -> &[u8] {
//         self.cipher.as_bytes()
//     }
//
//     fn open(&self) -> repo::crypto::cipher::CipherData {
//         if let None = self.nonce {
//             panic!("alskfjklsdjf");
//         }
//
//         let nonce = &self.nonce.as_ref().unwrap();
//         repo::crypto::cipher::CipherData {
//             nonce: hex::decode(nonce).unwrap().try_into().unwrap(),
//             ciphertext: hex::decode(&self.cipher).unwrap(),
//         }
//     }
//
//     fn nonce(&self) -> Option<Vec<u8>> {
//         match &self.nonce {
//             None => None,
//             Some(k) => {
//                 let nonce = hex::decode(k).unwrap();
//                 Some(nonce)
//             }
//         }
//     }
//
//     fn update(&mut self, d: repo::crypto::cipher::CipherData) {
//         self.cipher = hex::encode(d.ciphertext);
//         self.nonce = Some(hex::encode(d.nonce));
//     }
// }
//
// #[actix_web::get("/test")]
// async fn test_route() -> impl actix_web::Responder {
//     use repo::crypto::cipher;
//     let mut foo = Foo {
//         bar: "bar".to_owned(),
//         cipher: "baz".to_owned(),
//         nonce: None,
//     };
//
//     println!("before seal call");
//     dbg!(&foo);
//
//     cipher::Cipher::new(cipher::KeyType::RowKey)
//         .seal(&mut foo)
//         .unwrap();
//
//     println!("after seal call");
//     dbg!(&foo);
//
//     let decrypted = cipher::Cipher::new(cipher::KeyType::RowKey)
//         .open(&foo)
//         .unwrap();
//
//     println!("after open");
//     dbg!(&foo);
//     dbg!(std::str::from_utf8(&decrypted).unwrap());
//
//     return std::str::from_utf8(&foo.bar.as_ref()).unwrap().to_owned();
// }

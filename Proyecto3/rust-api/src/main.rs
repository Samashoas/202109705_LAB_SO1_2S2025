use actix_web::{web, App, HttpServer, Responder};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct WeatherTweet {
    municipality: String,
    temperature: i32,
    humidity: i32,
    weather: String,
}

async fn create_tweet(tweet: web::Json<WeatherTweet>) -> impl Responder {
    println!("Received tweet: {:?}", tweet);
    web::Json(serde_json::json!({ "status": "success" }))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server at http://127.0.0.1:8080");
    
    HttpServer::new(|| {
        App::new()
            .route("/tweets", web::post().to(create_tweet)) // Ruta POST para recibir tweets
    })
    .bind("127.0.0.1:8080")? // Puerto en el que escuchará la API
    .run()
    .await
}

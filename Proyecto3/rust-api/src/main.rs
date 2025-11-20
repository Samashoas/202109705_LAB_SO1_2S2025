use actix_web::{web, App, HttpServer, Result, HttpResponse, middleware::Logger};
use serde::{Deserialize, Serialize};
use reqwest;

#[derive(Serialize, Deserialize, Debug)]
struct WeatherTweet {
    municipality: String,
    temperature: i32,
    humidity: i32,
    weather: String,
}

#[derive(Serialize)]
struct Response {
    message: String,
}

async fn health() -> Result<HttpResponse> {
    Ok(HttpResponse::Ok().json(Response {
        message: "Rust API is running".to_string(),
    }))
}

async fn receive_weather(tweet: web::Json<WeatherTweet>) -> Result<HttpResponse> {
    println!("Received tweet: {:?}", tweet);
    
    // Enviar al deployment de Go
    let client = reqwest::Client::new();
    let go_service_url = std::env::var("GO_SERVICE_URL")
        .unwrap_or_else(|_| "http://go-service:8081".to_string());
    
    match client
        .post(&format!("{}/weather", go_service_url))
        .json(&tweet.into_inner())
        .send()
        .await
    {
        Ok(response) => {
            if response.status().is_success() {
                Ok(HttpResponse::Ok().json(Response {
                    message: "Weather data processed successfully".to_string(),
                }))
            } else {
                Ok(HttpResponse::InternalServerError().json(Response {
                    message: "Failed to process weather data".to_string(),
                }))
            }
        }
        Err(e) => {
            eprintln!("Error sending to Go service: {}", e);
            Ok(HttpResponse::InternalServerError().json(Response {
                message: "Internal server error".to_string(),
            }))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    
    println!("Starting Rust API server at http://0.0.0.0:8080");
    
    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .route("/health", web::get().to(health))
            .route("/weather", web::post().to(receive_weather))
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
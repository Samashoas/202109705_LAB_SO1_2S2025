package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "go-weather-api/proto"
)

type WeatherTweet struct {
    Municipality string `json:"municipality"`
    Temperature  int    `json:"temperature"`
    Humidity     int    `json:"humidity"`
    Weather      string `json:"weather"`
}

type Response struct {
    Message string `json:"message"`
}

func publishToKafkaGRPC(tweet WeatherTweet) error {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return err
    }
    defer conn.Close()
    client := pb.NewWeatherServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err = client.ProcessWeather(ctx, &pb.WeatherRequest{
        Municipality: tweet.Municipality,
        Temperature:  int32(tweet.Temperature),
        Humidity:     int32(tweet.Humidity),
        Weather:      tweet.Weather,
    })
    return err
}

func publishToRabbitGRPC(tweet WeatherTweet) error {
    conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return err
    }
    defer conn.Close()
    client := pb.NewWeatherServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err = client.ProcessWeather(ctx, &pb.WeatherRequest{
        Municipality: tweet.Municipality,
        Temperature:  int32(tweet.Temperature),
        Humidity:     int32(tweet.Humidity),
        Weather:      tweet.Weather,
    })
    return err
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Message: "Go API is running"})
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
    var tweet WeatherTweet

    if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    log.Printf("Received weather data: %+v", tweet)

    // Solo llamadas gRPC
    if err := publishToKafkaGRPC(tweet); err != nil {
        log.Printf("Error publishing to Kafka via gRPC: %v", err)
    }
    if err := publishToRabbitGRPC(tweet); err != nil {
        log.Printf("Error publishing to RabbitMQ via gRPC: %v", err)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Message: "Weather data processed via gRPC"})
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler).Methods("GET")
    r.HandleFunc("/weather", weatherHandler).Methods("POST")

    fmt.Printf("Go API server starting on port %s\n", port)
    fmt.Println("Configured as gRPC client for Kafka and RabbitMQ writers")
    log.Fatal(http.ListenAndServe(":"+port, r))
}
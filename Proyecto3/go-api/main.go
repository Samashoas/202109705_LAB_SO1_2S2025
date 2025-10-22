package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
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
    
    // TODO: Aquí implementaremos gRPC, Kafka y RabbitMQ
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Message: "Weather data received"})
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
    log.Fatal(http.ListenAndServe(":"+port, r))
}
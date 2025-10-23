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
    "github.com/segmentio/kafka-go"
    "github.com/streadway/amqp"
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

type WeatherProcessor struct {
    kafkaWriter  *kafka.Writer
    rabbitConn   *amqp.Connection
    rabbitCh     *amqp.Channel
    grpcClient   pb.WeatherServiceClient
}

func NewWeatherProcessor() (*WeatherProcessor, error) {
    wp := &WeatherProcessor{}
    
    // Configurar Kafka
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }
    
    wp.kafkaWriter = &kafka.Writer{
        Addr:     kafka.TCP(kafkaBroker),
        Topic:    "weather-tweets",
        Balancer: &kafka.LeastBytes{},
    }
    
    // Configurar RabbitMQ
    rabbitURL := os.Getenv("RABBITMQ_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }
    
    var err error
    wp.rabbitConn, err = amqp.Dial(rabbitURL)
    if err != nil {
        log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
    } else {
        wp.rabbitCh, err = wp.rabbitConn.Channel()
        if err != nil {
            log.Printf("Warning: Failed to open RabbitMQ channel: %v", err)
        } else {
            // Declarar cola
            _, err = wp.rabbitCh.QueueDeclare(
                "weather-queue", // name
                true,           // durable
                false,          // delete when unused
                false,          // exclusive
                false,          // no-wait
                nil,            // arguments
            )
            if err != nil {
                log.Printf("Warning: Failed to declare RabbitMQ queue: %v", err)
            }
        }
    }
    
    // Configurar gRPC client (opcional para ahora)
    grpcAddr := os.Getenv("GRPC_SERVICE_ADDR")
    if grpcAddr != "" {
        conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            log.Printf("Warning: Failed to connect to gRPC service: %v", err)
        } else {
            wp.grpcClient = pb.NewWeatherServiceClient(conn)
        }
    }
    
    return wp, nil
}

func (wp *WeatherProcessor) PublishToKafka(ctx context.Context, tweet WeatherTweet) error {
    messageBytes, err := json.Marshal(tweet)
    if err != nil {
        return err
    }
    
    return wp.kafkaWriter.WriteMessages(ctx,
        kafka.Message{
            Key:   []byte(tweet.Municipality),
            Value: messageBytes,
        },
    )
}

func (wp *WeatherProcessor) PublishToRabbitMQ(tweet WeatherTweet) error {
    if wp.rabbitCh == nil {
        return fmt.Errorf("RabbitMQ channel not available")
    }
    
    messageBytes, err := json.Marshal(tweet)
    if err != nil {
        return err
    }
    
    return wp.rabbitCh.Publish(
        "",              // exchange
        "weather-queue", // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        messageBytes,
            Timestamp:   time.Now(),
        },
    )
}

func (wp *WeatherProcessor) ProcessWeatherData(tweet WeatherTweet) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Publicar en Kafka
    if err := wp.PublishToKafka(ctx, tweet); err != nil {
        log.Printf("Error publishing to Kafka: %v", err)
    } else {
        log.Printf("Published to Kafka: %s", tweet.Municipality)
    }
    
    // Publicar en RabbitMQ
    if err := wp.PublishToRabbitMQ(tweet); err != nil {
        log.Printf("Error publishing to RabbitMQ: %v", err)
    } else {
        log.Printf("Published to RabbitMQ: %s", tweet.Municipality)
    }
    
    // Llamar gRPC si está disponible
    if wp.grpcClient != nil {
        req := &pb.WeatherRequest{
            Municipality: tweet.Municipality,
            Temperature:  int32(tweet.Temperature),
            Humidity:     int32(tweet.Humidity),
            Weather:      tweet.Weather,
        }
        
        resp, err := wp.grpcClient.ProcessWeather(ctx, req)
        if err != nil {
            log.Printf("Error calling gRPC service: %v", err)
        } else {
            log.Printf("gRPC response: %s", resp.Message)
        }
    }
    
    return nil
}

func (wp *WeatherProcessor) Close() {
    if wp.kafkaWriter != nil {
        wp.kafkaWriter.Close()
    }
    if wp.rabbitCh != nil {
        wp.rabbitCh.Close()
    }
    if wp.rabbitConn != nil {
        wp.rabbitConn.Close()
    }
}

var processor *WeatherProcessor

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
    
    // Procesar datos (Kafka, RabbitMQ, gRPC)
    if err := processor.ProcessWeatherData(tweet); err != nil {
        log.Printf("Error processing weather data: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Message: "Weather data processed successfully"})
}

func main() {
    var err error
    processor, err = NewWeatherProcessor()
    if err != nil {
        log.Fatalf("Failed to initialize weather processor: %v", err)
    }
    defer processor.Close()
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    
    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler).Methods("GET")
    r.HandleFunc("/weather", weatherHandler).Methods("POST")
    
    fmt.Printf("Go API server starting on port %s\n", port)
    fmt.Println("Configured with Kafka, RabbitMQ, and gRPC client support")
    log.Fatal(http.ListenAndServe(":"+port, r))
}
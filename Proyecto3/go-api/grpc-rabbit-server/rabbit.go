package main

import (
    "context"
    "log"
    "net"
    "os"
    "encoding/json"

    "google.golang.org/grpc"
    pb "go-weather-api/proto"
    "github.com/streadway/amqp"
)

type server struct {
    pb.UnimplementedWeatherServiceServer
    rabbitCh *amqp.Channel
}

func (s *server) ProcessWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
    msg := map[string]interface{}{
        "municipality": req.Municipality,
        "temperature":  req.Temperature,
        "humidity":     req.Humidity,
        "weather":      req.Weather,
    }
    body, _ := json.Marshal(msg)
    err := s.rabbitCh.Publish(
        "", "weather-queue", false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        log.Printf("Error publishing to RabbitMQ: %v", err)
        return &pb.WeatherResponse{Message: "Failed to publish", Success: false}, err
    }
    log.Printf("Published to RabbitMQ: %s", req.Municipality)
    return &pb.WeatherResponse{Message: "Published to RabbitMQ", Success: true}, nil
}

func main() {
    rabbitURL := os.Getenv("RABBITMQ_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }
    conn, err := amqp.Dial(rabbitURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open RabbitMQ channel: %v", err)
    }
    _, err = ch.QueueDeclare("weather-queue", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterWeatherServiceServer(grpcServer, &server{rabbitCh: ch})
    log.Println("gRPC RabbitMQ server listening on :50052")
    grpcServer.Serve(lis)
}
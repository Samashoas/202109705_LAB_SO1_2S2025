package main

import (
    "context"
    "log"
    "net"
    "os"

    "encoding/json"
    "google.golang.org/grpc"
    pb "go-weather-api/proto"
    "github.com/segmentio/kafka-go"
)

type server struct {
    pb.UnimplementedWeatherServiceServer
    kafkaWriter *kafka.Writer
}

func (s *server) ProcessWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
    msg := map[string]interface{}{
        "municipality": req.Municipality,
        "temperature":  req.Temperature,
        "humidity":     req.Humidity,
        "weather":      req.Weather,
    }
    messageBytes, err := json.Marshal(msg)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        return &pb.WeatherResponse{Message: "Failed to marshal", Success: false}, err
    }
    err = s.kafkaWriter.WriteMessages(ctx,
        kafka.Message{
            Key:   []byte(req.Municipality),
            Value: messageBytes,
        },
    )
    if err != nil {
        log.Printf("Error publishing to Kafka: %v", err)
        return &pb.WeatherResponse{Message: "Failed to publish", Success: false}, err
    }
    log.Printf("Published to Kafka: %s", req.Municipality)
    return &pb.WeatherResponse{Message: "Published to Kafka", Success: true}, nil
}

func main() {
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }
    writer := &kafka.Writer{
        Addr:     kafka.TCP(kafkaBroker),
        Topic:    "weather-tweets",
        Balancer: &kafka.LeastBytes{},
    }
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterWeatherServiceServer(grpcServer, &server{kafkaWriter: writer})
    log.Println("gRPC server listening on :50051")
    grpcServer.Serve(lis)
}
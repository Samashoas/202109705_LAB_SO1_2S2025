package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/segmentio/kafka-go"
)

func main() {
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{kafkaBroker},
        Topic:     "weather-tweets",
        GroupID:   "weather-consumer-group",
        Partition: 0,
        MinBytes:  10e3, // 10KB
        MaxBytes:  10e6, // 10MB
    })

    fmt.Println("Kafka consumer started. Waiting for messages...")
    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message: %v", err)
            continue
        }
        fmt.Printf("Consumed from Kafka: %s\n", string(m.Value))
    }
}
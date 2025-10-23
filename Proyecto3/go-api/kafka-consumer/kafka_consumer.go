package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/segmentio/kafka-go"
    "github.com/redis/go-redis/v9"
)

func main() {
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    ctx := context.Background()

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{kafkaBroker},
        Topic:     "weather-tweets",
        GroupID:   "weather-consumer-group",
        Partition: 0,
        MinBytes:  10e3,
        MaxBytes:  10e6,
    })

    fmt.Println("Kafka consumer started. Waiting for messages...")
    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message: %v", err)
            continue
        }
        fmt.Printf("Consumed from Kafka: %s\n", string(m.Value))
        // Guarda en Valkey
        err = rdb.RPush(ctx, "weather_kafka", string(m.Value)).Err()
        if err != nil {
            log.Printf("Error storing in Valkey: %v", err)
        }
        // Expira la lista en 6 minutos (360 segundos)
        err = rdb.Expire(ctx, "weather_kafka", 360).Err()
        if err != nil {
            log.Printf("Error setting expiration in Valkey: %v", err)
        }
        log.Printf("Stored message in Valkey: %s", string(m.Value))
    }
}
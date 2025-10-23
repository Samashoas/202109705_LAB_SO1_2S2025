package main

import (
    "fmt"
    "log"
    "os"
    "github.com/streadway/amqp"
    "github.com/redis/go-redis/v9"
    "context"
)

func main() {
    rabbitURL := os.Getenv("RABBITMQ_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }

    valkeyAddr := os.Getenv("VALKEY_ADDR")
    if valkeyAddr == "" {
        valkeyAddr = "localhost:6379"
    }
    rdb := redis.NewClient(&redis.Options{
        Addr: valkeyAddr,
    })
    ctx := context.Background()

    conn, err := amqp.Dial(rabbitURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    msgs, err := ch.Consume(
        "weather-queue", "", true, false, false, false, nil,
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    fmt.Println("RabbitMQ consumer started. Waiting for messages...")
    for msg := range msgs {
        fmt.Printf("Consumed from RabbitMQ: %s\n", string(msg.Body))
        // Guarda en Valkey
        err := rdb.RPush(ctx, "weather_rabbitmq", string(msg.Body)).Err()
        if err != nil {
            log.Printf("Error storing in Valkey: %v", err)
        }
        // Expira la lista en 360 segundos
        err = rdb.Expire(ctx, "weather_rabbitmq", 3600).Err()
        if err != nil {
            log.Printf("Error setting expiration in Valkey: %v", err)
        }
        log.Printf("Stored message in Valkey: %s", string(msg.Body))
    }
}
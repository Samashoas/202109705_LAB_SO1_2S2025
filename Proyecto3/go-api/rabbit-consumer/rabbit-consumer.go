package main

import (
    "fmt"
    "log"
    "os"
    "github.com/streadway/amqp"
)

func main() {
    rabbitURL := os.Getenv("RABBITMQ_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }

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
        "weather-queue", // queue
        "",              // consumer
        true,            // auto-ack
        false,           // exclusive
        false,           // no-local
        false,           // no-wait
        nil,             // args
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    fmt.Println("RabbitMQ consumer started. Waiting for messages...")
    for msg := range msgs {
        fmt.Printf("Consumed from RabbitMQ: %s\n", string(msg.Body))
    }
}
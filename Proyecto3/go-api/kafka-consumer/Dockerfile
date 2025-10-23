FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
COPY ./proto ./proto
RUN go build -o kafka-consumer kafka_consumer.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/kafka-consumer .
COPY --from=builder /app/proto ./proto
CMD ["./kafka-consumer"]
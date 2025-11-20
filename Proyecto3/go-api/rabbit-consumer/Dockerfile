FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
COPY ./proto ./proto
RUN go build -o rabbit-consumer rabbit-consumer.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rabbit-consumer .
COPY --from=builder /app/proto ./proto
CMD ["./rabbit-consumer"]
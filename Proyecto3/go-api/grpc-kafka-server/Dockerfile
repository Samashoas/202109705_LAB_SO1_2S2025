FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
COPY ./proto ./proto
RUN go mod download
RUN go build -o kafka-server kafka.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/kafka-server .
COPY --from=builder /app/proto ./proto
EXPOSE 50051
CMD ["./kafka-server"]
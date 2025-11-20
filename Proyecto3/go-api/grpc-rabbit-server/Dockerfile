FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
COPY ./proto ./proto
RUN go build -o rabbit-server rabbit.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rabbit-server .
COPY --from=builder /app/proto ./proto
EXPOSE 50052
CMD ["./rabbit-server"]
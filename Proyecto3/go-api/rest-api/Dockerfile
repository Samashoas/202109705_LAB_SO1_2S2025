FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
COPY ./proto ./proto
RUN go mod download
RUN go build -o rest-api main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rest-api .
COPY --from=builder /app/proto ./proto
EXPOSE 8081
CMD ["./rest-api"]
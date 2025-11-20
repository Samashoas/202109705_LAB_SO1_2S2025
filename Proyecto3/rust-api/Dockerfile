FROM rust:1.90-slim-bullseye as builder
WORKDIR /app

# Instala las dependencias necesarias para compilar OpenSSL crates
RUN apt-get update && apt-get install -y pkg-config libssl-dev

COPY . .
RUN cargo build --release

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/target/release/rust-api .
EXPOSE 8080
CMD ["./rust-api"]
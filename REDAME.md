Servicio de Deploy en GO
cd go-api
go run main.go

Servicio de API en RUST
cd rust-api
GO_SERVICE_URL=http://localhost:8081 cargo run

Servicio de Locust
source locust-env/bin/activate
locust -f locust_tweets.py --host=http://localhost:8080
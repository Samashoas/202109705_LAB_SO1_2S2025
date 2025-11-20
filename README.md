Servicio de Deploy en GO
cd go-api
go run main.go

Servicio de API en RUST
cd rust-api
GO_SERVICE_URL=http://localhost:8081 cargo run

Servicio de Locust
source locust-env/bin/activate
locust -f locust_tweets.py --host=http://localhost:8080

cp -r ../proto ./proto
cp ../go.mod ../go.sum 

docker build -t rabbit-consummer:latest .

curl http://136.114.218.161:5000/v2/_catalog

docker tag rust-api:latest 136.114.218.161:5000/rust-api:latest
docker push 136.114.218.161:5000/rust-api:latest

docker tag go-client:latest 136.114.218.161:5000/go-client:latest
docker push 136.114.218.161:5000/go-client:latest

docker tag grpc-kafka:latest 136.114.218.161:5000/grpc-kafka:latest
docker push 136.114.218.161:5000/grpc-kafka:latest

docker tag grpc-rabbit:latest 136.114.218.161:5000/grpc-rabbit:latest
docker push 136.114.218.161:5000/grpc-rabbit:latest

docker tag kafka-consummer:latest 136.114.218.161:5000/kafka-consummer:latest
docker push 136.114.218.161:5000/kafka-consummer:latest

docker tag rabbit-consummer:latest 136.114.218.161:5000/rabbit-consummer:latest
docker push 136.114.218.161:5000/rabbit-consummer:latest

docker build -t kafka-consummer:latest .


docker run -e GO_SERVICE_URL=http://192.168.1.126:8081 -p 8080:8080 rust-api:latest

docker run -p 8081:8081 \
  -e KAFKA_GRPC_ADDR=192.168.1.126:50051 \
  -e RABBIT_GRPC_ADDR=192.168.1.126:50052 \
  go-client:latest

docker run -d --name grpc-kafka-server \
  -e KAFKA_BROKER=192.168.1.126:9092 \
  -p 50051:50051 \
  grpc-kafka:latest

docker run -d --name grpc-rabbit-server \
  -e RABBITMQ_URL=amqp://guest:guest@192.168.1.126:5672/ \
  -p 50052:50052 \
  grpc-rabbit:latest

docker run -d --name kafka-consumer \
  -e KAFKA_BROKER=192.168.1.126:9092 \
  -e VALKEY_ADDR=192.168.1.126:6379 \
  kafka-consumer:latest

docker run -d --name rabbit-consumer \
  -e RABBITMQ_URL=amqp://guest:guest@192.168.1.126:5672/ \
  -e VALKEY_ADDR=192.168.1.126:6379 \
  rabbit-consumer:latest


1. Configura la consulta en el panel
Selecciona tu data source Redis.
Elige el tipo de consulta compatible (por ejemplo, List → LRANGE, si está disponible).
Ingresa la clave de tu lista, por ejemplo:
Start: 0
Stop: -1
Haz clic en Run.
2. Aplica transformaciones
Ve a la pestaña Transformations.
Usa Parse JSON para extraer los campos (municipality, humidity, etc.).
Aplica Group By para agrupar por municipality.
Aplica Aggregate para calcular el promedio de humidity por municipio.
3. Configura la visualización
Selecciona el panel de Bar Chart.
Eje X: municipality
Eje Y: promedio de humidity
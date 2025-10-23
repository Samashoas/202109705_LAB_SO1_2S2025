# Manual Técnico — Proyecto 3: Sistema Distribuido de Datos Meteorológicos

## 1. Arquitectura General

El sistema está compuesto por microservicios desarrollados en Go y Rust, que se comunican mediante gRPC, Kafka y RabbitMQ. Los datos procesados se almacenan en Valkey (compatible con Redis) y se visualizan en Grafana. El despliegue se realiza mediante Docker, Zot Registry y Kubernetes.

### Componentes Principales
- **REST API (Go):** Recibe datos meteorológicos y los envía a los escritores gRPC.
- **gRPC Writers (Go):** Publican datos en Kafka y RabbitMQ.
- **Kafka/RabbitMQ Consumers (Go):** Procesan mensajes y almacenan en Valkey con expiración.
- **API Principal (Rust):** Expone endpoints y lógica adicional.
- **Valkey:** Almacén de datos con expiración automática.
- **Grafana:** Visualización de datos usando el plugin Redis Data Source.
- **Zot Registry:** Registro OCI para imágenes Docker.
- **Kubernetes:** Orquestación y despliegue en la nube.

## 2. Estructura de Carpetas
```
Proyecto3/
  go-api/
    rest-api/
    grpc-kafka-server/
    grpc-rabbit-server/
    kafka-consumer/
    rabbit-consumer/
    proto/
  rust-api/
  locust/
```

## 3. Configuración de Entorno

Todos los servicios usan variables de entorno para direcciones y credenciales:
- `KAFKA_ADDR`, `RABBIT_ADDR`, `VALKEY_ADDR`, `GRPC_KAFKA_ADDR`, `GRPC_RABBIT_ADDR`, etc.
- Ejemplo en Docker Compose:
```yaml
  environment:
    - KAFKA_ADDR=kafka:9092
    - VALKEY_ADDR=valkey:6379
```

## 4. Compilación y Construcción de Imágenes

### Go
```bash
cd Proyecto3/go-api/rest-api
GOOS=linux GOARCH=amd64 go build -o rest-api
```

### Rust
```bash
cd Proyecto3/rust-api
cargo build --release
```

### Docker
```bash
# Construir imagen Go
cd Proyecto3/go-api/rest-api
sudo docker build -t zot.local:5000/rest-api:latest .

# Construir imagen Rust
cd Proyecto3/rust-api
sudo docker build -t zot.local:5000/rust-api:latest .
```

## 5. Publicación en Zot Registry

Configura Docker para permitir HTTP (inseguro):
```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "insecure-registries": ["zot.local:5000"]
}
EOF
sudo systemctl restart docker
```

Publica la imagen:
```bash
sudo docker push zot.local:5000/rest-api:latest
```

## 6. Despliegue Local con Docker Compose

Ejemplo para brokers y Valkey:
```yaml
version: '3.8'
services:
  kafka:
    image: bitnami/kafka:latest
    ports:
      - "9092:9092"
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
  valkey:
    image: valkey/valkey:latest
    ports:
      - "6379:6379"
```

## 7. Despliegue en Kubernetes

Ejemplo de Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rest-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rest-api
  template:
    metadata:
      labels:
        app: rest-api
    spec:
      containers:
      - name: rest-api
        image: zot.local:5000/rest-api:latest
        env:
        - name: GRPC_KAFKA_ADDR
          value: "grpc-kafka-server:50051"
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
```

## 8. Persistencia y Expiración en Valkey

Los consumidores almacenan datos usando RPush y configuran expiración:
```go
client.RPush(ctx, "weather", value)
client.Expire(ctx, "weather", 3600) // 1 hora
```

## 9. Visualización en Grafana

1. Instala el plugin Redis Data Source.
2. Configura la fuente apuntando a Valkey (`valkey:6379`).
3. Importa el dashboard desde `grafana/dashboards/system-monitor.json`.

## 10. Solución de Problemas

- **Errores de compilación Go:** Verifica `go.mod` y dependencias.
- **Errores Rust:** Asegura edición 2021 y dependencias correctas (`openssl`, `pkg-config`).
- **Problemas Docker:** Verifica contexto y variables de entorno.
- **Push a Zot falla:** Revisa configuración de `insecure-registries`.
- **Conexión Valkey/Grafana:** Verifica puertos y credenciales.

## 11. Buenas Prácticas

- Usa variables de entorno para configuración.
- Mantén los proto files en la carpeta `proto/` y cópialos en el contexto de build.
- Define recursos en Kubernetes para evitar sobrecarga.
- Documenta endpoints y flujos de datos.

## 12. Referencias
- [Valkey](https://valkey.io/)
- [Grafana Redis Data Source](https://grafana.com/grafana/plugins/redis-datasource/)
- [Zot Registry](https://zotregistry.io/)
- [Go](https://go.dev/)
- [Rust](https://www.rust-lang.org/)
- [Kafka](https://kafka.apache.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Kubernetes](https://kubernetes.io/)

---

# toll-calculator

### Kafka
```
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true bitnami/kafka:latest
```


### protobuf
```
brew install protoc-gen-go
```

### gRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Prometheus
```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
```

### Prometheus golang client
```
go get github.com/prometheus/client_golang/prometheus
```
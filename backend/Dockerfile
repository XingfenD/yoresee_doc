FROM golang:1.24-alpine AS dev
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache wget && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    wget -O /usr/local/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.22/grpc_health_probe-linux-amd64 && \
    chmod +x /usr/local/bin/grpc_health_probe

COPY go.mod ./
RUN go mod download

EXPOSE 8080
EXPOSE 9090
EXPOSE 2345

CMD ["sh", "-c", "go mod download && go run ./cmd/migrate && go run ./cmd/db_init && go run ./cmd/es_init && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient --continue --log ./cmd/main.go"]

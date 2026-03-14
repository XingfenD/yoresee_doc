FROM golang:1.24-alpine AS dev
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080
EXPOSE 9090
EXPOSE 2345

CMD ["sh", "-c", "go mod download && go run ./cmd/migrate && go run ./cmd/db_init && dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient --continue --log ./cmd/main.go"]

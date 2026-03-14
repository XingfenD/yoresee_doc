FROM golang:1.24-alpine AS dev
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

EXPOSE 1234

CMD ["sh", "-c", "go mod download && go run main.go"]

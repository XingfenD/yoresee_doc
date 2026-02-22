FROM golang:1.24-alpine AS dev
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

EXPOSE 8080
EXPOSE 2345

CMD ["sh", "-c", "go mod download && go install github.com/go-delve/delve/cmd/dlv@latest && go run ./cmd/migrate && dlv debug --headless --listen=:2345 --api-version=2 ./cmd/main.go"]

FROM golang:1.24-alpine AS builder

WORKDIR /app

# config the go proxy
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
COPY . .
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -p 1 -o main.out cmd/main.go && \
    CGO_ENABLED=0 go build -p 1 -o migrate.out ./cmd/migrate

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main.out .
COPY --from=builder /app/migrate.out .

EXPOSE 8080

CMD ["sh", "-c", "./migrate.out && ./main.out"]
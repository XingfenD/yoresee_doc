FROM golang:1.24-alpine AS dev
WORKDIR /app

EXPOSE 8080

CMD ["sh", "-c", "go mod download && go run ./cmd/migrate && go run cmd/main.go"]

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
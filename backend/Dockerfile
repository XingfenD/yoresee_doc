FROM golang:1.24-alpine AS builder

WORKDIR /app

# config the go proxy
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -p 1 -v -o main.out cmd/main.go && \
    CGO_ENABLED=0 go build -p 1 -v -o create_admin.out cmd/create_admin/main.go && \
    CGO_ENABLED=0 go build -p 1 -v -o migrate.out cmd/migrate/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main.out .
COPY --from=builder /app/create_admin.out .
COPY --from=builder /app/migrate.out .
COPY config.yaml .

EXPOSE 8080

CMD ["sh", "-c", "./migrate.out && ./create_admin.out && ./main.out"]

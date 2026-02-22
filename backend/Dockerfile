FROM golang:1.24-alpine AS dev
WORKDIR /app

EXPOSE 8080

CMD ["sh", "-c", "go mod download && go run cmd/migrate/main.go && go run cmd/init_config/main.go && go run cmd/init_permission/main.go && go run cmd/create_admin/main.go && go run cmd/main.go"]

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
    CGO_ENABLED=0 go build -p 1 -o create_admin.out cmd/create_admin/main.go && \
    CGO_ENABLED=0 go build -p 1 -o init_permission.out cmd/init_permission/main.go && \
    CGO_ENABLED=0 go build -p 1 -o init_config.out cmd/init_config/main.go && \
    CGO_ENABLED=0 go build -p 1 -o migrate.out cmd/migrate/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main.out .
COPY --from=builder /app/create_admin.out .
COPY --from=builder /app/init_permission.out .
COPY --from=builder /app/init_config.out .
COPY --from=builder /app/migrate.out .

EXPOSE 8080

CMD ["sh", "-c", "./migrate.out && ./init_config.out && ./init_permission.out && ./create_admin.out && ./main.out"]
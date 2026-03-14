FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o collab-go main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/collab-go .

EXPOSE 1235
CMD ["./collab-go"]

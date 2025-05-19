FROM golang:1.21-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go build -o weather-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/weather-service .

COPY local.yaml ./local.yaml

ENV CONFIG_PATH=/app/local.yaml

CMD ["./weather-service"]

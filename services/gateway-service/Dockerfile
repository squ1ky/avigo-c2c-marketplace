FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/gateway-service

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/gateway-service .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./gateway-service"]
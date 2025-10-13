FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/gateway-service

EXPOSE 8080

CMD ["./gateway-service"]
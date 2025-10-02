FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/user-service

EXPOSE 8081

CMD ["./user-service"]
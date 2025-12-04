FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -tags dynamic -o /app/listing-service ./cmd/listing-service/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates librdkafka wget

WORKDIR /app

COPY --from=builder /app/listing-service .

COPY --from=builder /app/internal/db/postgres/migrations ./migrations

COPY .env* ./


EXPOSE 8083

CMD ["./listing-service"]

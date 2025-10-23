FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o /build/notification-service \
    ./cmd/notification/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/notification-service .

COPY --from=builder /app/internal/migrations ./migrations

COPY --from=builder /app/assets ./assets

EXPOSE 8081

CMD ["./notification-service"]

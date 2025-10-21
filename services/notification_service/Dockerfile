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



# FROM golang:1.25.1-alpine AS builder

# RUN adduser -D -g '' appuser

# WORKDIR /build
# COPY go.mod go.sum ./
# RUN go mod download
# RUN go mod verify 

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -ldflags="-w -s" \
#     -a \
#     -installsuffix cgo \
#     -o /build/notification-service \
#     ./cmd/server/main.go

# FROM alpine:3.19

# RUN apk --no-cache add ca-certificates tzdata

# COPY --from=builder /etc/passwd /etc/passwd
# COPY --from=builder /etc/group/ /etc/group

# WORKDIR /app

# COPY --from=builder /build/migrations ./migrations
# COPY --from=builder /build/assets ./assets

# # Устанавливаем права на файлы
# RUN chown -R appuser:appuser /app

# # Переключаемся на непривилегированного пользователя
# USER appuser

# # Health check (можно добавить HTTP эндпоинт для проверки)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD pgrep -f notification-service || exit 1

# # Запускаем приложение
# ENTRYPOINT ["./notification-service"]
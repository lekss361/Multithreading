# ---- STAGE 1: builder ----
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Копируем только main.go и билдим без модулей
ENV GO111MODULE=off
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o myapp main.go

# ---- STAGE 2: runtime ----
FROM scratch

# Копируем собранный бинарник
COPY --from=builder /app/myapp /myapp

# Запуск бинарника при старте контейнера
ENTRYPOINT ["/myapp"]

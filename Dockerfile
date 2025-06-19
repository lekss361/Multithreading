# ---- STAGE 1: build ----
FROM golang:1.21-alpine AS builder

# рабочая директория внутри контейнера
WORKDIR /app

# кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# копируем весь исходник и собираем бинарник
COPY . .
# собираем статически с отключённым cgo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# ---- STAGE 2: runtime ----
FROM scratch

# если нужно, можно использовать маленький alpine:
# FROM alpine:latest
# RUN apk add --no-cache ca-certificates

# копируем бинарник из предыдущей стадии
COPY --from=builder /app/myapp /myapp

# укажите, на каком порту ваше приложение слушает (по умолчанию 80)
# EXPOSE 80

# путь к исполняемому файлу
ENTRYPOINT ["/myapp"]

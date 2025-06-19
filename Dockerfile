# ---- STAGE 1: build ----
FROM golang:1.21-alpine AS builder
WORKDIR /app

# кешируем модули
COPY go.mod go.sum ./
RUN go mod download

# копируем исходники и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# ---- STAGE 2: runtime ----
FROM scratch
# либо, если нужно SSL/CA, используйте alpine:
# FROM alpine:latest
# RUN apk add --no-cache ca-certificates

COPY --from=builder /app/myapp /myapp
ENTRYPOINT ["/myapp"]

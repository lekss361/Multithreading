FROM golang:1.21-alpine AS builder
WORKDIR /app

- # кешируем модули
- COPY go.mod go.sum ./
- RUN go mod download

# копируем весь исходник сразу
COPY . .

# если у вас пока нет go.mod/go.sum, можно просто собрать
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

FROM scratch
COPY --from=builder /app/myapp /myapp
ENTRYPOINT ["/myapp"]

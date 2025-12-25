# Сборка
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/todo-server ./cmd/todo-server/main.go

# Запуск
FROM alpine
USER root

WORKDIR /home/app

COPY --from=builder /bin/todo-server ./
COPY .env.example .env

EXPOSE 8080

ENTRYPOINT ["./todo-server"]
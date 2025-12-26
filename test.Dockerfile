FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go test ./... -v

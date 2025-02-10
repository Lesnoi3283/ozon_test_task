# builder
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go


# runner
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .

COPY .env .env

RUN chmod +x main

CMD ["sh", "-c", "./main -m=${IN_MEMORY_STORAGE}"]

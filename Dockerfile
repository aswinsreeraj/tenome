FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o tenome ./cmd/server

FROM debian:bookworm-slim

RUN apt update && \
    apt install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=builder /app/tenome .

CMD ["./tenome"]
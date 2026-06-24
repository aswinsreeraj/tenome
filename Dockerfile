FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o tenome ./cmd/server/main.go

FROM debain:bookworm-slim

WORKDIR /app

COPY --from=builder /app/tenome .

CMD ["./tenome"]
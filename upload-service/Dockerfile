# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Shanghai

COPY --from=builder /server /app/server
COPY config.yaml.example /app/config.yaml

EXPOSE 50054

CMD ["./server", "-c", "./config.yaml"]

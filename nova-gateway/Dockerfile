# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /nova-gateway cmd/server/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install ca-certificates for TLS
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# Set timezone
ENV TZ=Asia/Shanghai

# Copy binary from builder
COPY --from=builder /nova-gateway /app/nova-gateway

# Copy config
COPY config.yaml.example /app/config.yaml

EXPOSE 8080

CMD ["./nova-gateway", "-c", "./config.yaml"]

# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /web cmd/server/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install ca-certificates for TLS
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# Set timezone
ENV TZ=Asia/Shanghai

# Copy binary from builder
COPY --from=builder /web /app/web

# Copy config and resources
COPY configs/ /app/configs/
COPY resources/ /app/resources/

EXPOSE 8081

CMD ["./web"]
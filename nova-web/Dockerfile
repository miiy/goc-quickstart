# Frontend build stage
FROM node:22-bookworm-slim AS frontend-builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY tsconfig.json vite.config.ts ./
COPY resources/static ./resources/static
COPY frontend ./frontend
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}
RUN npm run build

# Go build stage
FROM golang:1.26 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
COPY --from=frontend-builder /app/dist ./dist

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
COPY config.yaml.example /app/config.yaml
COPY resources/ /app/resources/
COPY --from=frontend-builder /app/dist /app/dist

EXPOSE 8081

CMD ["./web", "-c", "./config.yaml"]

#!/bin/bash

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_DIR="$ROOT_DIR/logs"
PID_FILE="$LOG_DIR/.pids"

mkdir -p "$LOG_DIR"

stop_services() {
    if [ -f "$PID_FILE" ]; then
        echo "Stopping all services..."
        while read pid; do
            # Kill the process and all its children
            children=$(pgrep -P "$pid" 2>/dev/null)
            kill "$pid" 2>/dev/null
            for child in $children; do
                # go run -> compiled binary, kill the real process
                grandchildren=$(pgrep -P "$child" 2>/dev/null)
                kill "$child" "$grandchildren" 2>/dev/null
            done
        done < "$PID_FILE"
        wait 2>/dev/null
        rm -f "$PID_FILE"
        echo "All services stopped."
    fi
}

cleanup() {
    stop_services
    exit 0
}
trap cleanup EXIT INT TERM

start_service() {
    local name=$1
    local dir=$2
    local cmd=$3

    echo "Starting $name..."
    (cd "$ROOT_DIR/$dir" && eval "$cmd") > "$LOG_DIR/${name}.log" 2>&1 &
    local pid=$!
    echo "$pid" >> "$PID_FILE"
    echo "  $name started (PID: $pid, log: $LOG_DIR/${name}.log)"
}

rm -f "$PID_FILE"

# gRPC microservices
start_service "nova-auth"   "nova-auth"   "go run ./cmd/server/main.go -c ./config.yaml"
start_service "nova-user"   "nova-user"   "go run ./cmd/server/main.go -c ./config.yaml"
start_service "nova-post"   "nova-post"   "go run ./cmd/server/main.go -c ./config.yaml"
start_service "nova-file"   "nova-file" "go run ./cmd/server/main.go -c ./config.yaml"

sleep 2

# HTTP gateway
start_service "nova-gateway"    "nova-gateway"    "go run ./cmd/server/main.go -c ./config.yaml"

# Web frontend
start_service "nova-web"         "nova-web"       "go run ./cmd/server/main.go -c ./config.yaml"

echo ""
echo "All services started. Logs: $LOG_DIR/"
echo "Press Ctrl+C to stop all services."

wait

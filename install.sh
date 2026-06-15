#!/bin/bash

set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_PID_FILE="$PROJECT_DIR/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_DIR/.frontend.pid"
BACKEND_LOG="$PROJECT_DIR/.backend.log"
FRONTEND_LOG="$PROJECT_DIR/.frontend.log"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; }

is_running() {
    local pid_file="$1"
    if [ -f "$pid_file" ]; then
        local pid
        pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
    fi
    return 1
}

stop_process() {
    local name="$1"
    local pid_file="$2"
    if is_running "$pid_file"; then
        local pid
        pid=$(cat "$pid_file")
        log "Stopping $name (PID $pid)..."
        kill "$pid" 2>/dev/null || true
        sleep 1
        kill -9 "$pid" 2>/dev/null || true
        rm -f "$pid_file"
    fi
    # Also kill any process holding port 8000 (backend) to be safe
    if [ "$name" = "backend" ]; then
        lsof -ti:8000 2>/dev/null | xargs kill -9 2>/dev/null || true
    fi
}

start_backend() {
    if is_running "$BACKEND_PID_FILE"; then
        local pid
        pid=$(cat "$BACKEND_PID_FILE")
        warn "Backend already running (PID $pid), skipping."
        return
    fi

    log "Building backend..."
    cd "$PROJECT_DIR"
    go build -o .backend_bin . 2>&1 | tee -a "$BACKEND_LOG"

    log "Starting backend..."
    ./.backend_bin >> "$BACKEND_LOG" 2>&1 &
    echo $! > "$BACKEND_PID_FILE"
    log "Backend started (PID $(cat "$BACKEND_PID_FILE")), log: $BACKEND_LOG"
}

start_frontend() {
    if is_running "$FRONTEND_PID_FILE"; then
        local pid
        pid=$(cat "$FRONTEND_PID_FILE")
        warn "Frontend already running (PID $pid), skipping."
        return
    fi

    log "Installing frontend dependencies..."
    cd "$PROJECT_DIR/h5"
    npm install --silent

    log "Starting frontend dev server..."
    npm run dev >> "$FRONTEND_LOG" 2>&1 &
    echo $! > "$FRONTEND_PID_FILE"
    log "Frontend started (PID $(cat "$FRONTEND_PID_FILE")), log: $FRONTEND_LOG"
}

wait_for_frontend() {
    local url="http://localhost:5173"
    log "Waiting for frontend to be ready at $url ..."
    local i=0
    while ! curl -s "$url" > /dev/null 2>&1; do
        sleep 1
        i=$((i+1))
        if [ $i -ge 30 ]; then
            err "Frontend did not start within 30 seconds."
            return 1
        fi
    done
    log "Frontend is ready."
}

open_browser() {
    local url="http://localhost:5173"
    log "Opening browser at $url ..."
    if command -v open &>/dev/null; then
        open "$url"
    elif command -v xdg-open &>/dev/null; then
        xdg-open "$url"
    else
        warn "Could not detect a way to open the browser. Please open $url manually."
    fi
}

case "${1:-start}" in
    start)
        log "=== Starting game services ==="
        start_backend
        start_frontend
        wait_for_frontend
        open_browser
        log "=== Game is running ==="
        ;;
    stop)
        log "=== Stopping game services ==="
        stop_process "backend" "$BACKEND_PID_FILE"
        stop_process "frontend" "$FRONTEND_PID_FILE"
        rm -f "$PROJECT_DIR/.backend_bin"
        log "=== Game stopped ==="
        ;;
    status)
        if is_running "$BACKEND_PID_FILE"; then
            log "Backend:  running (PID $(cat "$BACKEND_PID_FILE"))"
        else
            warn "Backend:  not running"
        fi
        if is_running "$FRONTEND_PID_FILE"; then
            log "Frontend: running (PID $(cat "$FRONTEND_PID_FILE"))"
        else
            warn "Frontend: not running"
        fi
        ;;
    restart)
        "$0" stop
        sleep 1
        "$0" start
        ;;
    *)
        echo "Usage: $0 {start|stop|status|restart}"
        exit 1
        ;;
esac

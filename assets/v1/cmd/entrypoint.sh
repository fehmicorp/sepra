
#!/bin/sh
set -eu

export MINIO_ROOT_USER="${APP_ACCESS_KEY_ID:-minioadmin}"
export MINIO_ROOT_PASSWORD="${APP_SECRET_ACCESS_KEY:-minioadmin}"

APP_PORT="${APP_PORT:-9098}"
ADMIN_PORT="${ADMIN_PORT:-9091}"
DATA_DIR="${APP_DATA_DIR:-/data}"
BACKUP_DIR="${APP_BACKUP_DIR:-/backup}"
CF_TOKEN="${CF_TOKEN:-}"

MINIO_PID=""
CF_PID=""

mkdir -p "$DATA_DIR" "$BACKUP_DIR"

cleanup() {
    echo "Stopping SEPRA services..."

    if [ -n "$CF_PID" ]; then
        kill "$CF_PID" 2>/dev/null || true
    fi

    if [ -n "$MINIO_PID" ]; then
        kill "$MINIO_PID" 2>/dev/null || true
    fi

    if [ -n "$CF_PID" ]; then
        wait "$CF_PID" 2>/dev/null || true
    fi

    if [ -n "$MINIO_PID" ]; then
        wait "$MINIO_PID" 2>/dev/null || true
    fi
}

trap cleanup EXIT
trap 'exit 0' INT TERM

echo "Starting MinIO..."
minio server "$DATA_DIR" \
    --address ":9090" \
    --console-address ":$ADMIN_PORT" &
MINIO_PID=$!

echo "Waiting for MinIO readiness..."

i=0
until wget -q -O /dev/null \
    http://127.0.0.1:9090/minio/health/ready
do
    if ! kill -0 "$MINIO_PID" 2>/dev/null; then
        echo "ERROR: MinIO exited during startup."
        exit 1
    fi

    i=$((i + 1))
    if [ "$i" -ge 60 ]; then
        echo "ERROR: MinIO readiness timeout."
        exit 1
    fi

    sleep 1
done

echo "MinIO is ready."

# Do not copy raw MinIO data files as a backup.
# Use an S3-aware backup process for production.

echo "Initializing SEPRA bucket..."
/app/main

echo "SEPRA initialization completed."

if [ -n "$CF_TOKEN" ]; then
    echo "Starting Cloudflare Tunnel..."

    /usr/local/bin/cloudflared tunnel \
        --no-autoupdate \
        --protocol http2 \
        run --token "$CF_TOKEN" &
    CF_PID=$!

    echo "Cloudflare Tunnel process started."
else
    echo "CF_TOKEN not set; running without Cloudflare Tunnel."
fi

echo "Initializing HTTP server..."
/app/http

# Monitor both services. Exit if either required process dies.
while true; do
    if ! kill -0 "$MINIO_PID" 2>/dev/null; then
        echo "ERROR: MinIO stopped."
        exit 1
    fi

    if [ -n "$CF_PID" ] &&
       ! kill -0 "$CF_PID" 2>/dev/null; then
        echo "ERROR: Cloudflare Tunnel stopped."
        exit 1
    fi

    sleep 2
done

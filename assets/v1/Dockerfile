FROM alpine:3.20
WORKDIR /app
# Install required certificates and curl
RUN apk --no-cache add ca-certificates curl

# Download the MinIO binary from the official GitHub release assets
RUN curl -L -o /usr/bin/minio https://github.com/minio/minio/releases/latest/download/minio-linux-amd64 && \
    chmod +x /usr/bin/minio

# Define build arguments with defaults
ARG MINIO_ROOT_USER=minioadmin
ARG MINIO_ROOT_PASSWORD=minioadmin
ARG USE_SSL=false
ARG MAX_CONNECTIONS=100
ARG DEBUG_MODE=false
ARG TIMEOUT=30s
ARG REGION=ap-south-1
ARG PORT=9090
ARG MGMTPORT=9091

# Expose runtime environment variables to the container
ENV MINIO_ROOT_USER=${MINIO_ROOT_USER}
ENV MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD}

# Expose API and Console ports
EXPOSE ${PORT}
EXPOSE ${MGMTPORT}

# Start the MinIO server using the specified port and console management port
ENTRYPOINT ["minio", "server", "/data", "--address", ":9090", "--console-address", ":9091"]
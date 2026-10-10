FROM golang:1.26.2-alpine AS minio-builder
WORKDIR /app
RUN apk add --no-cache git
COPY cmd/minio/go.mod cmd/minio/go.sum ./
RUN go mod download
COPY cmd/minio .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o minio .
FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=minio-builder /app/minio /app/minio
COPY cmd/minio/config.yaml /app/config.yaml
EXPOSE 9090
ENTRYPOINT ["./minio", "--config-dir", "/app", "server", "/data"]
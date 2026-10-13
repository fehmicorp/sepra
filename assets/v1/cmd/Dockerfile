FROM golang:1.26.4-alpine AS minio-builder
RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0
RUN go install github.com/minio/minio@latest
FROM golang:1.26.4-alpine AS app-builder
WORKDIR /app
COPY ./http/go.mod ./http/go.sum ./
RUN go mod download
COPY ./app/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .
FROM golang:1.26.4-alpine AS http-builder
WORKDIR /app
COPY ./http/go.mod ./http/go.sum ./
RUN go mod download
COPY ./http/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .
FROM cloudflare/cloudflared:latest AS cloudflared-builder
FROM alpine:3.22
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=minio-builder /go/bin/minio /usr/local/bin/minio
COPY --from=app-builder /out/app /app/main
COPY --from=http-builder /out/app /app/http
COPY --from=cloudflared-builder /usr/local/bin/cloudflared /usr/local/bin/cloudflared
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod 0755 \
    /usr/local/bin/minio \
    /usr/local/bin/cloudflared \
    /app/entrypoint.sh \
    && mkdir -p /data /backup \
    && sed -i 's/\r$//' /app/entrypoint.sh
EXPOSE 9090 9091
VOLUME ["/data", "/backup"]
ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]

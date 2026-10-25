ARG VERSION=v1

FROM golang:1.26.4-alpine AS minio-builder
RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0
RUN go install github.com/minio/minio@latest

FROM golang:1.26.4-alpine AS app-builder
ARG VERSION
WORKDIR /app
COPY ./${VERSION}/cmd/app/go.mod ./${VERSION}/cmd/app/go.sum ./
RUN go mod download
COPY ./${VERSION}/cmd/app/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM golang:1.26.4-alpine AS http-builder
ARG VERSION
WORKDIR /app
COPY ./${VERSION}/cmd/http/go.mod ./${VERSION}/cmd/http/go.sum ./
RUN go mod download
COPY ./${VERSION}/cmd/http/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM cloudflare/cloudflared:latest AS cloudflared-builder

FROM alpine:3.22
ARG VERSION
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=minio-builder /go/bin/minio /usr/local/bin/minio
COPY --from=app-builder /out/app /app/main
COPY --from=http-builder /out/app /app/http
COPY --from=cloudflared-builder /usr/local/bin/cloudflared /usr/local/bin/cloudflared
COPY ./${VERSION}/cmd/entrypoint.sh /app/entrypoint.sh
RUN chmod 0755 \
    /usr/local/bin/minio \
    /usr/local/bin/cloudflared \
    /app/entrypoint.sh \
    && mkdir -p /data /backup \
    && sed -i 's/\r$//' /app/entrypoint.sh

EXPOSE 9090 9091 9098
VOLUME ["/data", "/backup"]
ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
ARG VERSION=v1

FROM cloudflare/cloudflared:latest AS cloudflared-builder

FROM alpine:3.22
ARG VERSION
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=cloudflared-builder /usr/local/bin/cloudflared /usr/local/bin/cloudflared
COPY ./${VERSION}/cmd/entrypoint.sh /app/entrypoint.sh
RUN chmod 0755 \
    /usr/local/bin/cloudflared \
    /app/entrypoint.sh \
    && mkdir -p /data /backup \
    && sed -i 's/\r$//' /app/entrypoint.sh
    
EXPOSE 9001
VOLUME ["/data", "/backup"]
ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
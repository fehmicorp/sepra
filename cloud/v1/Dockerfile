# --- Global Build Arguments ---
ARG TYPE=local
ARG PORT=8080
ARG HOST=0.0.0.0
ARG APPDIR=dist

# --- Frontend Builder ---
FROM node:20-alpine AS builder
ARG APP
ARG TYPE
WORKDIR /app
ENV OUTPUT=${TYPE}
ENV APP=${APP}
COPY webpack/${APP}/package.json webpack/${APP}/package-lock.json* ./
RUN npm install
COPY webpack/${APP}/ ./
RUN npm run build

# --- Go Server Builder ---
FROM golang:1.26.4-alpine AS server
ARG APP
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY cmd/${APP}/go.mod cmd/${APP}/go.sum* ./
RUN go mod download
COPY cmd/${APP}/ ./
COPY --from=builder /app/out ./dist/auth
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .

# --- Final Production Image ---
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache libc6-compat ca-certificates
COPY --from=server /app/server ./
COPY --from=server /app/dist ./dist
EXPOSE 8080
CMD ["./server"]
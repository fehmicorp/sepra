# --- Build Arguments ---
ARG TYPE=local
ARG PORT=8080
ARG HOST=0.0.0.0
ARG APPDIR=dist

# --- Frontend Auth Builder ---
FROM node:20-alpine AS auth
ARG APP=auth
ARG TYPE
WORKDIR /app
ENV OUTPUT=${TYPE}
ENV APP=${APP}
COPY webpack/${APP}/package.json webpack/${APP}/package-lock.json* ./
RUN npm install
COPY webpack/${APP}/ ./
RUN npm run build

# --- Frontend Home Builder ---
FROM node:20-alpine AS home
ARG APP=home
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
ARG APP=app
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY cmd/${APP}/go.mod cmd/${APP}/go.sum* ./
RUN go mod download
COPY cmd/${APP}/ ./
COPY --from=auth /app/out ./dist/auth
COPY --from=home /app/out ./dist/home
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .
FROM golang:1.26.4-alpine AS pkg
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
COPY pkg/config/go.mod pkg/config/go.sum* ./
RUN go mod download
COPY pkg/config/ ./
RUN mkdir -p ./pkg
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./pkg/config.so .

# --- Final Production Image ---
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache libc6-compat ca-certificates
COPY --from=server /app/server ./
COPY --from=pkg /app/pkg ./pkg
COPY --from=server /app/dist ./dist
EXPOSE 8080
CMD ["./server"]
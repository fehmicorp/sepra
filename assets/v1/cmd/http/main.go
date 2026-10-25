package main

import (
	"context"
	"fmt"
	"http/env"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

var minioClient *minio.Client
var cfg *Config
var (
	serviceToken = env.GetString("APP_SECRET_ACCESS_TOKEN", "")
	Port         = env.GetInt("APP_PORT", 9098)
	Ver          = env.GetString("VERSION", "v1")
	redisClient  *redis.Client
	ctx          = context.Background()
)

func main() {
	if os.Getenv("REDIS_HOST") != "" {
		InitRedis()
	}
	minioClient, _ = Client()
	mux := http.NewServeMux()
	ViewOnly(mux)
	CreateDownloadUrl(mux)
	Download(mux)
	RegisterBucketRoutes(mux)
	RegisterFileRoutes(mux)
	serverAddr := fmt.Sprintf(":%d", Port)
	fmt.Printf("Starting HTTP server on %s...\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

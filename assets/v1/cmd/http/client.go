package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"http/env"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
)

func Client() (*minio.Client, error) {
	minIOPort := env.GetInt("MINIO_PORT", 9090)
	minIOHost := fmt.Sprintf("localhost:%d", minIOPort)

	accessKeyId := env.GetString("APP_ACCESS_KEY_ID", "minioadmin")
	secretAccessKey := env.GetString("APP_SECRET_ACCESS_KEY", "minioadmin")
	// secretAccessToken := env.GetString("APP_SECRET_ACCESS_TOKEN", "")

	useSSL := env.GetBool("APP_USE_SSL", false)
	maxConnections := env.GetInt("APP_MAX_CONNECTIONS", 100)

	client, err := minio.New(minIOHost, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure:       useSSL,
		MaxRetries:   maxConnections,
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}
	return client, nil
}

func InitRedis() {
	redisHost := env.GetString("REDIS_HOST", "localhost:6379")
	redisPassword := env.GetString("REDIS_PASSWORD", "")
	redisDb := env.GetInt("REDIS_DB", 0)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       redisDb,
	})
}

func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/pkg/credentials"
	"github.com/minio/minio-go/v7"
)

var (
	minioClient     *minio.Client
	endpoint        = env.GetString("APP_ENDPOINT", "9090")
	accessKeyID     = env.GetString("APP_ACCESS_KEY_ID", "minioadmin")
	secretAccessKey = env.GetString("APP_SECRET_ACCESS_KEY", "minioadmin")
	useSSL          = env.GetBool("APP_USE_SSL", false)
	maxConnections  = env.GetInt("APP_MAX_CONNECTIONS", 10)
	debugMode       = env.GetBool("APP_DEBUG_MODE", false)
	timeout         = env.GetDuration("APP_TIMEOUT", 30*time.Second)
)

func main() {
	ctx := context.Background()
	// 1. Initialize MinIO client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Failed to initialize MinIO client:", err)
	}

	fmt.Println("Successfully connected to MinIO!")

	bucketName := "my-test-bucket"
	location := "us-east-1"
}

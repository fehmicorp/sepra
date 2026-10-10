package main

import (
	"context"
	"fmt"
	"log"
	"minio/env"
	"minio/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func main() {
	ctx := context.Background()

	// Load configuration
	rawCfg, err := env.LoadConfig("config.yaml", &Config{}, &Config{})
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg := rawCfg.(*Config)

	fmt.Printf("Config Loaded Successfully:\n")
	fmt.Printf("  Port:             %d\n", cfg.Port)
	fmt.Printf("  AccessKeyID:      %s\n", cfg.AccessKeyID)
	fmt.Printf("  MaxConnections:   %d\n", cfg.MaxConnections)
	fmt.Printf("  DebugMode:        %t\n", cfg.DebugMode)
	fmt.Printf("  Timeout:          %s\n", cfg.Timeout)

	// Ensure Data Directory
	dataDir, err := utils.EnsureDir(cfg.DataDir, true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Data directory: %s\n", dataDir)

	// Ensure Backup Directory
	backupDir, err := utils.EnsureDir(cfg.BackupDir, true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Backup directory: %s\n", backupDir)

	// Parse timeout duration to milliseconds
	minVal, err := utils.ConvertDuration(cfg.Timeout, "ms")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Timeout in milliseconds: %.0f ms\n", minVal)

	// Parse timeout as time.Duration for MinIO client options
	// parsedTimeout, err := time.ParseDuration(cfg.Timeout)
	// if err != nil {
	// 	parsedTimeout = 30 * time.Second // fallback
	// }

	// 1. Initialize MinIO client object (using v7 credentials)
	minioClientAddress := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	minioClient, err = minio.New(minioClientAddress, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure:       cfg.UseSSL,
		MaxRetries:   cfg.MaxConnections,
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		log.Fatalln("Failed to initialize MinIO client:", err)
	}

	// 2. Check and Create Bucket
	bucketName := "bin"
	location := cfg.Region
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalln("Error checking bucket existence:", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			log.Fatalf("Failed to create bucket %s: %v\n", bucketName, err)
		}
		fmt.Printf("Successfully created bucket: %s\n", bucketName)
	} else {
		fmt.Printf("Bucket %s already exists\n", bucketName)
	}
}

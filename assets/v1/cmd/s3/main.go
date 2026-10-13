package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"s3/env"
	"s3/utils"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

var cfg *Config

func main() {
	ctx := context.Background()
	// Load configuration
	rawCfg, err := env.LoadConfig("config1.yaml", &Config{}, &Config{})
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg = rawCfg.(*Config)

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

	// 1. Initialize MinIO client object (using v7 credentials)
	minioClient, err = minio.New("localhost:9090", &minio.Options{
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

	// Set up HTTP Multiplexer and register S3 proxy routes
	mux := http.NewServeMux()
	RegisterS3Routes(mux, cfg.SecretAccessToken)
	CreateDownloadUrl(mux, cfg.SecretAccessToken)
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Starting HTTP server on %s...\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func CreateDownloadUrl(mux *http.ServeMux, expectedToken string) {
	mux.HandleFunc("POST /{bucketName}/{dir...}", func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		token := r.Header.Get("X-SecretAccessToken")
		if token == "" || token != expectedToken {
			http.Error(w, "Unauthorized: Invalid or missing X-SecretAccessToken", http.StatusUnauthorized)
			return
		}
		// 2. Extract path parameters
		bucketName := r.PathValue("bucketName")
		objectName := r.PathValue("dir")
		if bucketName == "" || objectName == "" {
			http.Error(w, "Bad Request: Missing bucket name or file path", http.StatusBadRequest)
			return
		}

		// 3. Generate MinIO Presigned Download URL (valid for 2 hours)
		ctx := context.Background()
		reqParams := make(url.Values)
		presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, 2*time.Hour, reqParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate download URL: %v", err), http.StatusInternalServerError)
			return
		}

		// 4. Return the signed URL as a JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, presignedURL.String())))
	})
}

func RegisterS3Routes(mux *http.ServeMux, expectedToken string) {
	// Go 1.22+ standard router supports method specification (GET) and wildcards ({dir...})
	mux.HandleFunc("GET /{bucketName}/{dir...}", func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		token := r.Header.Get("X-SecretAccessToken")
		if token == "" || token != expectedToken {
			http.Error(w, "Unauthorized: Invalid or missing X-SecretAccessToken", http.StatusUnauthorized)
			return
		}

		// 2. Extract path parameters
		bucketName := r.PathValue("bucketName")
		objectName := r.PathValue("dir")

		if bucketName == "" || objectName == "" {
			http.Error(w, "Bad Request: Missing bucket name or file path", http.StatusBadRequest)
			return
		}

		// 3. Fetch object from MinIO
		ctx := context.Background()
		object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve object: %v", err), http.StatusInternalServerError)
			return
		}
		defer object.Close()

		// 4. Verify object metadata / existence
		objectInfo, err := object.Stat()
		if err != nil {
			http.Error(w, "Object not found", http.StatusNotFound)
			return
		}

		// 5. Set headers and stream response
		w.Header().Set("Content-Type", objectInfo.ContentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
		w.WriteHeader(http.StatusOK)

		if _, err := io.Copy(w, object); err != nil {
			// Handle client disconnection or network interruption during streaming
			return
		}
	})
}

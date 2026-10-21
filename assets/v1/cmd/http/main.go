package main

import (
	"context"
	"fmt"
	"http/env"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

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
	RegisterS3Routes(mux)
	CreateDownloadUrl(mux)
	serverAddr := fmt.Sprintf(":%d", Port)
	fmt.Printf("Starting HTTP server on %s...\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func CreateDownloadUrl(mux *http.ServeMux) {
	pattern := fmt.Sprintf("POST /%s/{bucketName}/{dir...}", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		token := r.Header.Get("X-SecretAccessToken")
		if token == "" || token != serviceToken {
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
		reqParams := make(url.Values)
		presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, 2*time.Hour, reqParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate download URL: %v", err), http.StatusInternalServerError)
			return
		}

		// 4. Generate a secure random unique download token
		downloadToken := generateSecureToken(32)

		// 5. Store the presigned URL in Redis with a 2-hour TTL
		redisKey := fmt.Sprintf("dl_token:%s", downloadToken)
		err = redisClient.Set(ctx, redisKey, presignedURL.String(), 2*time.Hour).Err()
		if err != nil {
			http.Error(w, "Failed to cache download token", http.StatusInternalServerError)
			return
		}

		// 6. Return the downloadable endpoint link as JSON
		downloadEndpoint := fmt.Sprintf("http://localhost:%d/%s/download?token=%s", Port, Ver, downloadToken)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"download_url":"%s"}`, downloadEndpoint)))

		// // 4. Return the signed URL as a JSON response
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// _, _ = w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, presignedURL.String())))
	})
}

func Download(mux *http.ServeMux, Ver string) {
	pattern := fmt.Sprintf("GET /%s/download", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract token from query parameters
		downloadToken := r.URL.Query().Get("token")
		if downloadToken == "" {
			http.Error(w, "Bad Request: Missing download token", http.StatusBadRequest)
			return
		}

		// 2. Verify token in Redis and retrieve the stored MinIO presigned URL
		redisKey := fmt.Sprintf("dl_token:%s", downloadToken)
		presignedURL, err := redisClient.Get(ctx, redisKey).Result()
		if err == redis.Nil {
			http.Error(w, "Unauthorized or Expired: Invalid download link", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// (Optional) Delete token immediately after use for single-use download links:
		// _ = redisClient.Del(ctx, redisKey)

		// 3. Redirect client to the presigned MinIO URL to initiate download/stream
		http.Redirect(w, r, presignedURL, http.StatusTemporaryRedirect)
	})
}

func RegisterS3Routes(mux *http.ServeMux) {
	pattern := fmt.Sprintf("GET /%s/{bucketName}/{dir...}", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		token := r.Header.Get("X-SecretAccessToken")
		if token == "" || token != serviceToken {
			http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
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

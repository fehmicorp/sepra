package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"http/env"

	"github.com/minio/minio-go/v7"
)

var minioClient *minio.Client
var cfg *Config
var (
	serviceToken = env.GetString("APP_SECRET_ACCESS_TOKEN", "")
	Port         = env.GetInt("APP_PORT", 9098)
)

func main() {
	fmt.Println(Port)
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
	mux.HandleFunc("POST /{bucketName}/{dir...}", func(w http.ResponseWriter, r *http.Request) {
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

func RegisterS3Routes(mux *http.ServeMux) {
	// Go 1.22+ standard router supports method specification (GET) and wildcards ({dir...})
	mux.HandleFunc("GET /{bucketName}/{dir...}", func(w http.ResponseWriter, r *http.Request) {
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

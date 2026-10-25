package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

func CreateDownloadUrl(mux *http.ServeMux) {
	pattern := fmt.Sprintf("POST /%s/{bucketName}/{dir...}", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		if !isAuthorizedRequest(r) {
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

func Download(mux *http.ServeMux) {
	pattern := fmt.Sprintf("GET /%s/download", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract token from query parameters
		if !isAuthorizedRequest(r) {
			http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
			return
		}
		downloadToken := r.URL.Query().Get("token")
		if downloadToken == "" {
			http.Error(w, "Bad Request: Missing download token", http.StatusBadRequest)
			return
		}

		// 2. Verify token in Redis and retrieve the stored bucket/object path
		redisKey := fmt.Sprintf("dl_token:%s", downloadToken)
		signedURL, err := redisClient.Get(ctx, redisKey).Result()
		if err == redis.Nil {
			http.Error(w, "Unauthorized or Expired: Invalid download link", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// (Optional) Uncomment below if you want the link to be strictly single-use:
		_ = redisClient.Del(ctx, redisKey)

		// 3. Fetch the file using the pre-signed MinIO URL
		resp, err := http.Get(signedURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to fetch object from storage backend", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// 4. Extract file name from the URL path or default to download
		u, err := url.Parse(signedURL)
		fileName := "downloaded-file"
		if err == nil {
			fileName = path.Base(u.Path)
		}

		// 5. Set headers to force file download and stream content
		if resp.Header.Get("Content-Type") != "" {
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		if resp.ContentLength > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		w.WriteHeader(http.StatusOK)

		// 6. Stream the response body directly to the client
		if _, err := io.Copy(w, resp.Body); err != nil {
			return
		}
	})
}

func ViewOnly(mux *http.ServeMux) {
	pattern := fmt.Sprintf("GET /%s/{bucketName}/{dir...}", Ver)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate X-SecretAccessToken header
		if !isAuthorizedRequest(r) {
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

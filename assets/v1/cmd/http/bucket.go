package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

func RegisterBucketRoutes(mux *http.ServeMux) {
	// 1. Create Bucket
	mux.HandleFunc(fmt.Sprintf("POST /%s/buckets/{bucketName}", Ver), createBucketHandler)

	// 2. Fetch / List Buckets
	mux.HandleFunc(fmt.Sprintf("GET /%s/buckets", Ver), listBucketsHandler)

	// 3. Delete Bucket (Supports ?mode=bin or ?mode=purge)
	mux.HandleFunc(fmt.Sprintf("DELETE /%s/buckets/{bucketName}", Ver), deleteBucketHandler)
}

// 1. Create Bucket
func createBucketHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	bucketName := r.PathValue("bucketName")
	if bucketName == "" {
		http.Error(w, "Bad Request: Missing bucket name", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil && exists {
		http.Error(w, "Conflict: Bucket already exists", http.StatusConflict)
		return
	}

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create bucket: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"Bucket '%s' created successfully"}`, bucketName)))
}

// 2. Fetch / List Buckets
func listBucketsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list buckets: %v", err), http.StatusInternalServerError)
		return
	}

	// Build a simple JSON array of bucket names and creation dates
	var sb strings.Builder
	sb.WriteString(`{"buckets":[`)
	for i, b := range buckets {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{"name":"%s","creation_date":"%s"}`, b.Name, b.CreationDate.Format(time.RFC3339)))
	}
	sb.WriteString(`]}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(sb.String()))
}

// 3. Delete Bucket (Move to Bin vs Purge Delete)
func deleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	bucketName := r.PathValue("bucketName")
	mode := r.URL.Query().Get("mode") // expected: "bin" or "purge"

	if bucketName == "" {
		http.Error(w, "Bad Request: Missing bucket name", http.StatusBadRequest)
		return
	}

	if mode == "" {
		mode = "bin" // Default to soft delete / move to bin for safety
	}

	ctx := context.Background()

	// Check if bucket exists
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil || !exists {
		http.Error(w, "Not Found: Bucket does not exist", http.StatusNotFound)
		return
	}

	trashBucket := ".trash-bin"

	switch mode {
	case "bin":
		// Move to Bin: Ensure trash bucket exists, copy all objects there, remove from source, then remove bucket
		_ = minioClient.MakeBucket(ctx, trashBucket, minio.MakeBucketOptions{})

		objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
		for obj := range objectCh {
			if obj.Err != nil {
				continue
			}
			// Copy object to trash with timestamp prefix
			src := minio.CopySrcOptions{Bucket: bucketName, Object: obj.Key}
			dst := minio.CopyDestOptions{Bucket: trashBucket, Object: fmt.Sprintf("%s/%s/%s", bucketName, time.Now().Format("20060102150405"), obj.Key)}
			_, _ = minioClient.CopyObject(ctx, dst, src)

			// Remove original object
			_ = minioClient.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		}

	case "purge":
		// Purge Delete: Recursively remove all objects in the bucket first
		objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
		for obj := range objectCh {
			if obj.Err != nil {
				continue
			}
			_ = minioClient.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		}

	default:
		http.Error(w, "Bad Request: Invalid delete mode. Use ?mode=bin or ?mode=purge", http.StatusBadRequest)
		return
	}

	// Finally, remove the empty bucket
	err = minioClient.RemoveBucket(ctx, bucketName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove bucket: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"Bucket '%s' successfully deleted using mode '%s'"}`, bucketName, mode)))
}

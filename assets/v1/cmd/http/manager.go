package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

func RegisterFileRoutes(mux *http.ServeMux) {
	// 1. Upload File to Bucket/Directory
	mux.HandleFunc(fmt.Sprintf("POST /%s/files/{bucketName}/{dir...}", Ver), uploadFileHandler)

	// 2. Make Directory (Creates a 0-byte marker object with trailing slash)
	mux.HandleFunc(fmt.Sprintf("POST /%s/files/mkdir/{bucketName}/{dir...}", Ver), makeDirHandler)

	// 3. Move File/Folder between buckets or paths
	mux.HandleFunc(fmt.Sprintf("POST /%s/files/move", Ver), moveFileHandler)

	// 4. Delete File/Folder (Moves to .trash-bin with restore mapping)
	mux.HandleFunc(fmt.Sprintf("DELETE /%s/files/{bucketName}/{dir...}", Ver), deleteFileHandler)

	// 5. Restore File from Trash back to original location
	mux.HandleFunc(fmt.Sprintf("POST /%s/files/restore", Ver), restoreFileHandler)
}

// 1. Upload File Handler
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	bucketName := r.PathValue("bucketName")
	objectName := r.PathValue("dir")
	if bucketName == "" || objectName == "" {
		http.Error(w, "Bad Request: Missing bucket name or destination path", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 32MB memory buffer)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Bad Request: Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request: Missing 'file' in form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := context.Background()
	// Upload file to MinIO
	_, err = minioClient.PutObject(ctx, bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"File '%s' uploaded successfully to bucket '%s'"}`, objectName, bucketName)))
}

// 2. Make Directory Handler
func makeDirHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	bucketName := r.PathValue("bucketName")
	dirPath := r.PathValue("dir")
	if bucketName == "" || dirPath == "" {
		http.Error(w, "Bad Request: Missing bucket name or directory path", http.StatusBadRequest)
		return
	}

	// Ensure directory path ends with a trailing slash
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	ctx := context.Background()
	// Create a 0-byte object to represent the directory structure
	_, err := minioClient.PutObject(ctx, bucketName, dirPath, strings.NewReader(""), 0, minio.PutObjectOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create directory: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"Directory '%s' created in bucket '%s'"}`, dirPath, bucketName)))
}

// 3. Move File Handler (JSON Body Request)
type MoveRequest struct {
	SrcBucket  string `json:"src_bucket"`
	SrcObject  string `json:"src_object"`
	DestBucket string `json:"dest_bucket"`
	DestObject string `json:"dest_object"`
}

func moveFileHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.SrcBucket == "" || req.SrcObject == "" || req.DestBucket == "" || req.DestObject == "" {
		http.Error(w, "Bad Request: Missing source or destination parameters", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Copy object to destination
	src := minio.CopySrcOptions{Bucket: req.SrcBucket, Object: req.SrcObject}
	dst := minio.CopyDestOptions{Bucket: req.DestBucket, Object: req.DestObject}
	_, err := minioClient.CopyObject(ctx, dst, src)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to copy object during move: %v", err), http.StatusInternalServerError)
		return
	}

	// Remove source object
	err = minioClient.RemoveObject(ctx, req.SrcBucket, req.SrcObject, minio.RemoveObjectOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Object copied but failed to delete source: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success","message":"File moved successfully"}`))
}

// 4. Delete File/Folder (Move to .trash-bin with Valkey Restore Point)
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	bucketName := r.PathValue("bucketName")
	objectPrefix := r.PathValue("dir")
	if bucketName == "" || objectPrefix == "" {
		http.Error(w, "Bad Request: Missing bucket name or file/folder path", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	trashBucket := ".trash-bin"
	_ = minioClient.MakeBucket(ctx, trashBucket, minio.MakeBucketOptions{})

	// List all objects matching the file or folder prefix
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: true,
	})

	movedCount := 0
	timestamp := time.Now().Format("20060102150405")

	for obj := range objectCh {
		if obj.Err != nil {
			continue
		}

		trashID := fmt.Sprintf("%s_%d_%s", timestamp, movedCount, generateSecureToken(6))
		trashObjectPath := fmt.Sprintf("trash/%s/%s", trashID, obj.Key)

		// Copy to trash bucket
		src := minio.CopySrcOptions{Bucket: bucketName, Object: obj.Key}
		dst := minio.CopyDestOptions{Bucket: trashBucket, Object: trashObjectPath}
		_, err := minioClient.CopyObject(ctx, dst, src)
		if err != nil {
			continue
		}

		// Store restore mapping in Valkey: trashID -> originalBucket|originalObjectPath
		restoreMapping := fmt.Sprintf("%s|%s", bucketName, obj.Key)
		_ = redisClient.Set(ctx, fmt.Sprintf("restore:%s", trashID), restoreMapping, 30*24*time.Hour).Err() // Retain for 30 days

		// Remove original object
		_ = minioClient.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		movedCount++
	}

	if movedCount == 0 {
		http.Error(w, "Not Found: No files found matching the specified path", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"Successfully moved %d item(s) to trash bin"}`, movedCount)))
}

// 5. Restore File Handler (JSON Body Request)
type RestoreRequest struct {
	TrashID string `json:"trash_id"`
}

func restoreFileHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized: Invalid Request", http.StatusUnauthorized)
		return
	}

	var req RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TrashID == "" {
		http.Error(w, "Bad Request: Missing or invalid trash_id", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	trashBucket := ".trash-bin"
	redisKey := fmt.Sprintf("restore:%s", req.TrashID)

	// Retrieve original location from Valkey
	originalLocation, err := redisClient.Get(ctx, redisKey).Result()
	if err != nil {
		http.Error(w, "Not Found: Invalid or expired trash ID", http.StatusNotFound)
		return
	}

	parts := strings.SplitN(originalLocation, "|", 2)
	if len(parts) != 2 {
		http.Error(w, "Internal Server Error: Malformed restore mapping", http.StatusInternalServerError)
		return
	}
	originalBucket := parts[0]
	originalObjectPath := parts[1]
	trashObjectPath := fmt.Sprintf("trash/%s/%s", req.TrashID, originalObjectPath)

	// Copy back to original bucket & directory path
	src := minio.CopySrcOptions{Bucket: trashBucket, Object: trashObjectPath}
	dst := minio.CopyDestOptions{Bucket: originalBucket, Object: originalObjectPath}
	_, err = minioClient.CopyObject(ctx, dst, src)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to restore object: %v", err), http.StatusInternalServerError)
		return
	}

	// Clean up trash object and Valkey mapping
	_ = minioClient.RemoveObject(ctx, trashBucket, trashObjectPath, minio.RemoveObjectOptions{})
	_ = redisClient.Del(ctx, redisKey)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"success","message":"File restored successfully to bucket '%s' at path '%s'"}`, originalBucket, originalObjectPath)))
}

package main

import (
	"fmt"
	"log"
	"minio/utils"
	"time"

	"github.com/minio/minio-go/v7"
)

var (
	minioClient *minio.Client
)

func main() {
	// rawCfg, err := env.LoadConfig("config.yaml", &Config{}, &Config{})
	// if err != nil {
	// 	log.Fatalf("Failed to load config: %v", err)
	// }

	// cfg := rawCfg.(*Config)
	// fmt.Printf("Config Loaded Successfully:\n")
	// fmt.Printf("  Port:             %d\n", cfg.Port)
	// fmt.Printf("  AccessKeyID:      %s\n", cfg.AccessKeyID)
	// fmt.Printf("  MaxConnections:   %d\n", cfg.MaxConnections)
	// fmt.Printf("  DebugMode:        %t\n", cfg.DebugMode)
	// fmt.Printf("  Timeout:          %s\n", cfg.Timeout)

	// minVal, err := utils.ConvertDuration(cfg.Timeout, "ms")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	// fmt.Printf("Timeout in milliseconds: %.0f ms\n", minVal)

	postTime, err := utils.CalcDuration(time.Time{}, "1", "d", false)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("After 1 Day:   %s\n", postTime.Format(time.RFC1123))

	// // 1. Initialize MinIO client object
	// minioClient, err := minio.New(cfg.Port, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	// 	Secure: cfg.UseSSL,
	// })
	// if err != nil {
	// 	log.Fatalln("Failed to initialize MinIO client:", err)
	// }

	// fmt.Println("Successfully connected to MinIO!")

	// bucketName := "my-test-bucket"
	// location := "us-east-1"
}

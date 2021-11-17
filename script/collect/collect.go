package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	bucketName      string
	accessKeyID     string
	secretAccessKey string

	minioClient *minio.Client
)

func init() {
	bucketName = os.Getenv("BUCKET_NAME")
	accessKeyID = os.Getenv("ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("SECRET_ACCESS_KEY")

	var err error
	minioClient, err = minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	resp, err := http.Get("http://127.0.0.1:8081/files")
	if err != nil {
		log.Fatal("http request error: ", err)
	}

	var body struct {
		Status int      `json:"status,omitempty"`
		Data   []string `json:"data,omitempty"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&body); err != nil {
		log.Fatal("http decode error: ", err)
	}

	ctx := context.Background()
	if exists, err := minioClient.BucketExists(ctx, bucketName); err == nil && !exists {
		if err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully created %s\n", bucketName)
	} else if err != nil {
		log.Fatal("check bucket error: ", err)
	}

	for _, file := range body.Data {
		fileResp, err := http.Get("http://127.0.0.1:8080/data/" + file)
		if err != nil {
			log.Println("error: ", err)
			continue
		}

		info, err := minioClient.PutObject(ctx, bucketName, file, fileResp.Body, fileResp.ContentLength, minio.PutObjectOptions{})
		if err != nil {
			log.Println("error: ", err)
			continue
		}

		log.Printf("Successfully uploaded %s of size %d\n", file, info.Size)
	}
}

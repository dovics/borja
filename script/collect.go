package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type LightData struct {
	Time  time.Time `json:"time,omitempty"`
	Value int       `json:"value,omitempty"`
}

var (
	bucketName      = "light-data"
	fileName        = "./light.db"
	accessKeyID     = "default_access_key_id"
	secretAccessKey = "default_secret_access_key"
)

func init() {
	bucketName = os.Getenv("BUCKET_NAME")
	fileName = os.Getenv("FILE_NAME")
	accessKeyID = os.Getenv("ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
}

func main() {
	var body struct {
		Status int          `json:"status,omitempty"`
		Data   []*LightData `json:"data,omitempty"`
	}
	resp, err := http.Get("http://127.0.0.1:8080/data")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(fileName); err != nil {
		log.Println(err)
	}

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTable := `
	CREATE TABLE light (
		time DATETIME NOT NULL PRIMARY KEY,
		  val INT NOT NULL
	);
	`

	db.Exec(createTable)
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(`INSERT INTO light(time, val) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}

	for _, elem := range body.Data {
		stmt.Exec(elem.Time, elem.Value)
	}

	objectName := time.Now().Format("20060102-150405")
	if err := UploadFile(bucketName, fileName, objectName); err != nil {
		log.Fatal(err)
	}
}

func UploadFile(bucketName string, filePath string, objectName string) error {
	minioClient, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	ctx := context.TODO()
	if exists, err := minioClient.BucketExists(ctx, bucketName); err == nil && !exists {
		if err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("Successfully created %s\n", bucketName)
	} else {
		return err
	}

	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

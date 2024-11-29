package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewMongo(ctx context.Context) *mongo.Client {
	fmt.Println("MONGO_URI", os.Getenv("MONGO_URI"))
	// MONGO_URI="mongodb://admin:admin@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@10.10.100.10:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func NewMinio(ctx context.Context) *minio.Client {
	/// 1. Initialize MinIO client

	endpoint := os.Getenv("MINIO_ENDPOINT") // Replace with your MinIO server
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

	fmt.Println("MINIO_ENDPOINT", endpoint)
	// useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
		panic(err)
	}
	return minioClient
}

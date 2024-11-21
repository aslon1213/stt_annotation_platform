package initializers

import (
	"context"
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
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

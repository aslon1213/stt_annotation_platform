package handlers

import (
	"context"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	users  *mongo.Collection
	audios *mongo.Collection
	jobs   *mongo.Collection
	minio  *minio.Client
	ctx    context.Context
}

func NewHandlers(db *mongo.Database, minio *minio.Client) *Handlers {

	// stt_works collection
	jobs := db.Collection("jobs")
	users := db.Collection("users")
	audios := db.Collection("audios")
	context := context.Background()

	return &Handlers{
		users:  users,
		audios: audios,
		jobs:   jobs,
		ctx:    context,
		minio:  minio,
	}
}

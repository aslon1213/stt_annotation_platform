package main

import (
	"context"
	"stt_work/handlers"
	"stt_work/initializers"
)

func main() {

	ctx := context.Background()
	initializers.LoadEnvs()

	client := initializers.NewMongo(ctx).Database("stt_works")
	minioClient := initializers.NewMinio(ctx)
	username := "admin"
	password := "admin"
	hls := handlers.NewHandlers(client, minioClient)

	hls.CreateUser(username, password)
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stt_work/initializers"
	"stt_work/models"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
)

func STTRequest(fileID string, minio_client *minio.Client) string {
	url := "http://10.10.100.10:8003"

	// read from minio
	obj, err := minio_client.GetObject(context.Background(), "agro-call-center", fileID, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
	defer obj.Close()
	// read to buffer
	buff := bytes.NewBuffer(nil)
	_, err = io.Copy(buff, obj)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", buff)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body := new(bytes.Buffer)
	output := map[string]interface{}{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body.Bytes(), &output)
	return output["text"].(string)
}

func main() {
	ctx := context.Background()
	initializers.LoadEnvs()
	minio_client := initializers.NewMinio(ctx)
	mongo_client := initializers.NewMongo(ctx)

	// hls := handlers.NewHandlers(mongo_client.Database("stt_works"), minio_client)

	// get all jobs
	jons := []models.Job{}
	cur, err := mongo_client.Database("stt_works").Collection("jobs").Find(ctx, bson.M{
		"initial_stt_done": false,
	})
	if err != nil {
		panic(err)
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &jons)
	if err != nil {
		panic(err)
	}
	for _, job := range jons {
		fmt.Println("Processing job: ", job.ID)
		transcript := STTRequest(job.AudioID.Hex(), minio_client)
		_, err := mongo_client.Database("stt_works").Collection("jobs").UpdateOne(ctx, bson.M{
			"_id": job.ID,
		}, bson.M{
			"$set": bson.M{
				"initial_stt_done": true,
				"stt_transcript":   transcript,
			},
		})
		if err != nil {
			panic(err)
		}
		// time.Sleep(1000000 * time.Second)
	}

}

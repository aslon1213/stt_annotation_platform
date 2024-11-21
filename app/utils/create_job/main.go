package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"stt_work/handlers"
	"stt_work/initializers"
	"stt_work/models"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func process_stt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buffer := bytes.Buffer{}

	// Copy the file data to buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return "", fmt.Errorf("error copying file data to buffer: %w", err)
	}

	resp, err := http.Post("http://localhost:8000/", "application/octet-stream", &buffer)
	if err != nil {
		return "", fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	buffer.Reset()
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error copying response body to buffer: %w", err)
	}

	output := map[string]interface{}{}
	err = json.Unmarshal(buffer.Bytes(), &output)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	text, ok := output["text"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: missing 'text'")
	}
	return text, nil
}

type Files struct {
	Files [][]string `json:"files"`
}

func worker(jobs <-chan []string, wg *sync.WaitGroup, hls *handlers.Handlers) {
	defer wg.Done()
	for file := range jobs {
		audio := &models.Audio{
			ID:       primitive.NewObjectID(),
			FilePath: file[1],
			Name:     file[0],
		}

		// fmt.Printf("Processing STT for file: %s\n", audio.FilePath)
		// stt_text, err := process_stt(audio.FilePath)
		// if err != nil {
		// 	log.Printf("Error processing STT for %s: %v\n", audio.FilePath, err)
		// 	continue
		// }

		audio.StorageUrl = fmt.Sprintf("http://%s/agro-call-center/%s", os.Getenv("MINIO_ACCESS_URL"), audio.ID.Hex())
		job := &models.Job{
			ID:              primitive.NewObjectID(),
			AudioID:         audio.ID,
			STTProcessed:    false,
			STTtranscript:   "",
			HumanProcessed:  false,
			HumanTranscript: "",
			WorkerID:        primitive.NilObjectID,
		}

		fmt.Printf("Creating job for audio: %s\n", audio.Name)
		if err := hls.CreateJob(job, audio); err != nil {
			log.Printf("Error creating job for %s: %v\n", audio.Name, err)
		}
	}
}

func main() {
	initializers.LoadEnvs()
	client := initializers.NewMongo(context.Background())
	minioClient := initializers.NewMinio(context.Background())
	hls := handlers.NewHandlers(client.Database("stt_works"), minioClient)

	// Read files list
	file, err := os.Open("files.json")
	if err != nil {
		log.Fatalf("Error opening files.json: %v", err)
	}
	defer file.Close()

	files := &Files{}
	buffer := bytes.Buffer{}
	io.Copy(&buffer, file)
	err = json.Unmarshal(buffer.Bytes(), files)
	if err != nil {
		log.Fatalf("Error unmarshalling files.json: %v", err)
	}

	// Create job channel and worker pool
	jobs := make(chan []string, len(files.Files))
	var wg sync.WaitGroup
	numWorkers := 4 // Number of workers (adjust as needed)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, hls)
	}

	// Send jobs to workers
	for _, file := range files.Files {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All jobs processed.")
}

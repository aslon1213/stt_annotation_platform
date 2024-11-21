package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"stt_work/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// uploadObject uploads a file to a specified bucket in MinIO
func uploadObject(client *minio.Client, bucketName, objectName, filePath, contentType string) error {
	ctx := context.Background()

	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Get file info to determine size
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	// Upload the file
	_, err = client.PutObject(ctx, bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (h *Handlers) UploadAudio(audio *models.Audio) error {
	bucket_name := "agro-call-center"
	object_name := audio.ID.Hex()
	file_path := audio.FilePath
	// object_url := fmt.Sprintf("http://localhost:9000/%s/%s", bucket_name, file_path)
	content_type := "audio/wav"
	err := uploadObject(h.minio, bucket_name, object_name, file_path, content_type)
	if err != nil {
		return err
	}

	// create a job object

	h.jobs.InsertOne(h.ctx, audio)
	return nil
}

func (h *Handlers) GetAudio(c *fiber.Ctx) error {
	bucket_name := "agro-call-center"
	fmt.Println("bucket_name: ", bucket_name)

	object_id := c.Params("id")

	obj, err := h.minio.GetObject(h.ctx, bucket_name, object_id, minio.GetObjectOptions{})
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer obj.Close()

	// Set the appropriate headers
	c.Set("Content-Type", "audio/wav")                                              // Set the content type
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", object_id)) // Optional: set filename

	// Stream the object to the response
	if _, err := io.Copy(c.Response().BodyWriter(), obj); err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return err // Added return statement to exit the function
	}

	return nil
}

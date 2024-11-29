package handlers

import (
	"fmt"
	"os"
	"strings"
	"stt_work/models"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handlers) GetJob(c *fiber.Ctx) error {
	return c.SendString("GetJob")
}

func (h *Handlers) GetJobs(c *fiber.Ctx) error {

	cursor, err := h.jobs.Find(h.ctx, bson.M{
		"human_processed": true,
	}, options.Find().SetLimit(50))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	jobs := []models.Job{}
	if err := cursor.All(h.ctx, &jobs); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("jobs", fiber.Map{
		"title": "Jobs",
		"jobs":  jobs,
	})

}

func (h *Handlers) JobDone(c *fiber.Ctx) error {

	job_id := c.Params("id")
	transcription := c.FormValue("transcription")
	fmt.Println("transcription: ", transcription)
	fmt.Println("job_id: ", job_id)

	transcription = strings.ReplaceAll(transcription, "\n", " ")

	user := c.Locals("user").(*models.User)
	// type Job struct {
	// 	ID              primitive.ObjectID `mongo:"_id" json:"id"`
	// 	AudioID         primitive.ObjectID `mongo:"audio_id" json:"audio_id"`
	// 	STTProcessed    bool               `mongo:"initial_stt_done" json:"initial_stt_done"`
	// 	STTtranscript   string             `mongo:"stt_transcript" json:"stt_transcript"`
	// 	HumanProcessed  bool               `mongo:"human_processed" json:"human_processed"`
	// 	HumanTranscript string             `mongo:"human_transcript" json:"human_transcript"`
	// 	WorkerID        primitive.ObjectID `mongo:"worker_username" json:"worker_username"`
	// }
	id, err := primitive.ObjectIDFromHex(job_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, err = h.jobs.UpdateByID(h.ctx, id, bson.M{
		"$set": bson.M{
			"human_processed":  true,
			"human_transcript": transcription,
			"worker_id":        user.ID,
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Redirect("/")
}

func (h *Handlers) UploadToMinio(filepath string, id primitive.ObjectID) error {
	// read the file

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// upload the file to minio
	// Get file info to determine size
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	info, err := h.minio.PutObject(h.ctx, "agro-call-center", id.Hex(), file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: "audio/wav",
	})
	if err != nil {
		return err
	}
	fmt.Println("info: ", info)
	return nil
}

func (h *Handlers) CreateJob(job *models.Job, audio *models.Audio) error {
	// insert audio first

	_, err := h.audios.InsertOne(h.ctx, audio)
	if err != nil {
		return err
	}

	_, err = h.jobs.InsertOne(h.ctx, job)
	fmt.Println("Inserted Job:")

	if err != nil {
		return err
	}

	// create an audio object

	h.UploadToMinio(audio.FilePath, audio.ID)

	return nil

}

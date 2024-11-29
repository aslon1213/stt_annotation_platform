package handlers

import (
	"stt_work/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handlers) Metrics(c *fiber.Ctx) error {

	// jobs which are done

	jobs_processed := []models.Job{}

	cur, err := h.jobs.Find(h.ctx, bson.M{"human_processed": true})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer cur.Close(h.ctx)
	err = cur.All(h.ctx, &jobs_processed)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	number_of_processed_jobs := len(jobs_processed)
	user_to_jobs_done := map[primitive.ObjectID]int{}
	for _, job := range jobs_processed {
		if _, ok := user_to_jobs_done[job.WorkerID]; ok {
			user_to_jobs_done[job.WorkerID] += 1
		} else {
			user_to_jobs_done[job.WorkerID] = 1
		}
	}
	users_name_to_jobs_done := map[string]int{}
	for user_id, jobs_done := range user_to_jobs_done {
		user := models.User{}
		err := h.users.FindOne(h.ctx, bson.M{"_id": user_id}).Decode(&user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		users_name_to_jobs_done[user.Username] = jobs_done
	}

	return c.JSON(fiber.Map{
		"number_of_processed_jobs": number_of_processed_jobs,
		"users_name_to_jobs_done":  users_name_to_jobs_done,
	},
	)
}

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
	jobs_unprocessed := []models.Job{}
	all_jobs := []models.Job{}

	cur, err := h.jobs.Find(h.ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error 1": err.Error(),
		})
	}
	defer cur.Close(h.ctx)
	err = cur.All(h.ctx, &all_jobs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error 2": err.Error(),
		})
	}

	for _, job := range all_jobs {
		if job.HumanProcessed {
			jobs_processed = append(jobs_processed, job)
		} else {
			jobs_unprocessed = append(jobs_unprocessed, job)
		}
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
	type UserNameToJob struct {
		Username string
		Jobsdone int
	}
	users_name_to_jobs_done := []UserNameToJob{}
	for user_id, jobs_done := range user_to_jobs_done {
		user := models.User{}
		err := h.users.FindOne(h.ctx, bson.M{
			"_id": user_id,
		}).Decode(&user)
		if err != nil {
			continue
		}
		users_name_to_jobs_done = append(users_name_to_jobs_done, UserNameToJob{
			Username: user.Username,
			Jobsdone: jobs_done,
		})
	}

	return c.Render("metrics", fiber.Map{
		"number_of_processed_jobs":   number_of_processed_jobs,
		"number_of_unprocessed_jobs": len(jobs_unprocessed),
		"all_jobs":                   len(all_jobs),
		"users_name_to_jobs_done":    users_name_to_jobs_done,
	})
}

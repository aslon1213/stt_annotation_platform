package main

import (
	"context"
	"stt_work/handlers"
	"stt_work/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v3"
)

func main() {

	// Create a new engine
	engine := django.New("views", ".django")
	app := fiber.New(
		fiber.Config{
			Views: engine,
		},
	)
	// Initialize default config
	app.Use(logger.New())

	// Or extend your config for customization
	// Logging remote IP and Port
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n", //
	}))
	// init mongo client
	// os.Getenv("MONGO_URI")
	ctx := context.Background()
	initializers.LoadEnvs()
	client := initializers.NewMongo(ctx)

	minioClient := initializers.NewMinio(ctx)
	hls := handlers.NewHandlers(client.Database("stt_works"), minioClient)

	app.Get("/metrics", hls.AuthMiddleware, hls.Metrics)

	jobs := app.Group("/jobs")
	users := app.Group("/users")
	audios := app.Group("/audios")

	jobs.Get("/", hls.AuthMiddleware, hls.GetJobs)

	//  get
	users.Post("/login", hls.Login)
	users.Get("/login", hls.LoginPage)
	app.Get("/", hls.AuthMiddleware, hls.IndexPage)

	//  get audio
	audios.Get("/:id", hls.AuthMiddleware, hls.GetAudio)
	jobs.Get("/:id", hls.AuthMiddleware, hls.GetJob)
	jobs.Get("/", hls.AuthMiddleware, hls.GetJobs)
	jobs.Post("/done/:id", hls.AuthMiddleware, hls.JobDone)

	app.Listen(":8080")
}

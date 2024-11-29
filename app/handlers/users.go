package handlers

import (
	"fmt"
	"os"
	"stt_work/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func (h *Handlers) Login(c *fiber.Ctx) error {

	// get user data
	username := c.FormValue("username")
	password := c.FormValue("password")
	fmt.Println("username: ", username)
	fmt.Println("password: ", password)
	user := models.User{}
	err := h.users.FindOne(h.ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid credentials 2",
		})
	}
	expire_time := time.Now().Add(24 * 5 * time.Hour) // 5 days
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "foo":  "bar",
		"user":   user.Username,
		"expire": expire_time,
	})

	// Sign and get the complete encoded token as a string using the secret
	fmt.Println("MY_SECRET: ", os.Getenv("MY_SECRET"))
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  expire_time,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	// generate a jwt token   ///////////////////////////////////////////////////////////////
	return c.Redirect("/")
}

func (h *Handlers) IndexPage(c *fiber.Ctx) error {

	jobs := []models.Job{}
	cursor, err := h.jobs.Find(h.ctx, bson.M{
		"human_processed": false,
	}, options.Find().SetLimit(50))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	defer cursor.Close(h.ctx)
	err = cursor.All(h.ctx, &jobs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Render("jobs", fiber.Map{
		"title": "Index",
		"jobs":  jobs,
	})
}

func (h *Handlers) CreateUser(username, password string) error {

	user := models.User{}
	user.Username = username

	user.ID = primitive.NewObjectID()
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(pass)

	_, err = h.users.InsertOne(h.ctx, user)
	if err != nil {
		return err
	}
	return nil
}

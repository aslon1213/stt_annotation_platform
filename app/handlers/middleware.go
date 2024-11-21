package handlers

import (
	"fmt"
	"os"
	"stt_work/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenValues struct {
	Username string
	Expire   int64
}

var connected_users_tokens = map[string]TokenValues{}

func (h *Handlers) AuthMiddleware(c *fiber.Ctx) error {

	tokenString := c.Cookies("Authorization", "")
	// fmt.Println("Cookies", string(c.Request().Header.Header()))
	if tokenString == "" {
		return c.Redirect("/users/login")
	}
	// fmt.Println("tokenString: ", tokenString)

	// check the token is valid or not
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		mysecret := os.Getenv("MY_SECRET")
		return []byte(mysecret), nil
	})

	if err != nil {

		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
			"err":   err.Error(),
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user := &models.User{}
		fmt.Println(
			"claims: ", claims,
		)
		username, ok := claims["user"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token",
			})

		}
		err := h.users.FindOne(h.ctx, bson.M{
			"username": username,
		}).Decode(user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		c.Locals("user", user)

		c.Next()

	} else {
		fmt.Println("claims: ", claims)
		c.Redirect("/users/login")
	}
	return nil
}

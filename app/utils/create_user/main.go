package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"stt_work/handlers"
	"stt_work/initializers"
	"fmt"
)

const passwordLength = 16 // Length of the generated password

// GenerateRandomPassword generates a random password of a specified length
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=<>?"
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[index.Int64()]
	}
	return string(password), nil
}

func main() {
	ctx := context.Background()
	initializers.LoadEnvs()

	client := initializers.NewMongo(ctx).Database("stt_works")
	minioClient := initializers.NewMinio(ctx)
	hls := handlers.NewHandlers(client, minioClient)


	for i := 2; i < 10; i++ {
		username := fmt.Sprintf("user%d",i)

		// Generate a random password
		password, err := GenerateRandomPassword(passwordLength)
		if err != nil {
			panic("Failed to generate a random password: " + err.Error())
		}

			
		hls.CreateUser(username, password)
	// Print the generated password (optional, for debugging or logging purposes)
	// Ensure to handle password logging securely in production!
		println("User", username ," ----- Generated password:", password)
	}
}


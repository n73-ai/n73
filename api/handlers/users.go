package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func VerifyAuthLink(c *fiber.Ctx) error {
	tokenString := c.Params("tokenString")

	token, err := utils.ParseAndValidateToken(tokenString, os.Getenv("SECRET_KEY"))
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(403).JSON(fiber.Map{
			"error": "Invalid token claim.",
		})
	}

	email := fmt.Sprint(claims["email"])
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing email claim.",
		})
	}

	exists, err := database.UserExists(email)

	secondsToByValid := 60 * 60 * 24 // 24 hours valid
	expDuration := time.Duration(secondsToByValid) * time.Second
	now := time.Now().UTC()
	exp := now.Add(expDuration).Unix()

	if exists {
		token, err := utils.GenerateJWT(email, secondsToByValid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"token": token,
			"exp":   exp,
			"email": email,
		})
	}

	userID := uuid.NewString()
	err = database.CreateUser(email, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	signToken, err := utils.GenerateJWT(email, secondsToByValid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token": signToken,
		"exp":   exp,
		"email": email,
	})
}

func AuthLink(c *fiber.Ctx) error {
	payload := struct {
		Email string `json:"email"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The email is required.",
		})
	}

	if len(payload.Email) > 55 {
		return c.Status(400).JSON(fiber.Map{
			"error": "The email should not have more than 55 characters.",
		})
	}

	secondsToByValid := 60 * 15 // 15 min valid

	token, err := utils.GenerateJWT(payload.Email, secondsToByValid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// token string, email string, subject string
	err = utils.SendEmail(token, payload.Email, "Welcome to AI Zustack")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

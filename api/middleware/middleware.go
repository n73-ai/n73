package middleware

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func User(c *fiber.Ctx) error {
	tokenString := utils.ExtractTokenFromHeader(c.Get("Authorization"))

	if tokenString == "" {
		return c.Status(401).SendString("You are not logged in.")
	}

	token, err := utils.ParseAndValidateToken(tokenString, os.Getenv("SECRET_KEY"))
	if err != nil {
		return c.Status(403).SendString(err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(403).SendString("Invalid token claim.")
	}

	user, err := database.GetUserByEmail(fmt.Sprint(claims["email"]))
	if err != nil {
		if err.Error() == "No user found with email "+fmt.Sprint(claims["email"]) {
			return c.Status(403).SendString("No user found with this token.")
		}
		return c.Status(500).SendString(err.Error())
	}

	c.Locals("scope", claims["scope"])
	c.Locals("user", &user)
	return c.Next()
}

func Admin(c *fiber.Ctx) error {
	tokenString := utils.ExtractTokenFromHeader(c.Get("Authorization"))

	if tokenString == "" {
		return c.Status(401).SendString("You are not logged in.")
	}

	token, err := utils.ParseAndValidateToken(tokenString, os.Getenv("SECRET_KEY"))
	if err != nil {
		return c.Status(403).SendString(err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(403).SendString("Invalid token claim.")
	}

	user, err := database.GetUserByEmail(fmt.Sprint(claims["email"]))
	if err != nil {
		if err.Error() == "No user found with email "+fmt.Sprint(claims["email"]) {
			return c.Status(403).SendString("No user found with this token.")
		}
		return c.Status(500).SendString(err.Error())
	}

	if user.Role != "admin" {
		return c.SendStatus(403)
	}

	c.Locals("scope", claims["scope"])
	c.Locals("user", &user)
	return c.Next()
}

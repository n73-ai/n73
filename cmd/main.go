package main

import (
	"ai-zustack/api"
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"log"
	"os"
)

func main() {
	required := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_NAME",
		"DB_PORT",
		"PORT",
		"EMAIL_SECRET_KEY",
		"ROOT_PATH",
		"SECRET_KEY",
		"ADMIN_JWT",
		"IP",
	}
	if err := utils.CheckRequiredEnv(required); err != nil {
		fmt.Printf("Environment error: %v.", err)
		os.Exit(1)
	}

	err := database.ConnectDB(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	if err != nil {
		fmt.Printf("DB connection error: %v.", err)
		os.Exit(1)
	}

	app := api.RunServer()
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	log.Fatal(app.Listen(port))
}

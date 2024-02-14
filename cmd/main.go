package main

import (
	"log"
    "github.com/joho/godotenv"

	"github.com/AlexLevus/telegram-bot/internal/app"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
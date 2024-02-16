package main

import (
	"log"
    "github.com/joho/godotenv"

	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
    "github.com/AlexLevus/telegram-bot/internal/controllers"
	"github.com/AlexLevus/telegram-bot/internal/routes"
	"github.com/AlexLevus/telegram-bot/internal/services"
	tele "gopkg.in/telebot.v3"
)

var (
	bot *tele.Bot

	BotController      controllers.BotController
	BotRouteController routes.BotRouteController

	chatService         services.ChatService
	chatCollection      *mongo.Collection

	pollCollection      *mongo.Collection
	suggestionCollection      *mongo.Collection
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

	mongoDbUri, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		log.Fatal("Добавьте в файл .env uri к базе MongoDB")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatal("Добавьте в файл .env название базы данных")
	}

	telegramApiToken, exists := os.LookupEnv("TELEGRAM_APITOKEN")
	if !exists {
		log.Fatal("Добавьте в файл .env api токен Телеграм")
	}

	clientOptions := options.Client().ApplyURI(mongoDbUri)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(dbName)
	chatCollection = database.Collection("chats")
	pollCollection = database.Collection("polls")
	suggestionCollection = database.Collection("suggestion")

	botSettings := tele.Settings{
		Token:  telegramApiToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(botSettings)
	if err != nil {
		log.Fatal(err)
	}

	bot = b

	chatService = services.NewChatService(chatCollection, ctx)

 
	BotController = controllers.NewBotController(bot, chatService)
	BotRouteController = routes.NewBotRouteController(BotController)
}

func main() {
	startTelegramBot()
}

func startTelegramBot() {
	BotRouteController.BotRoute(bot)

	bot.Start()
}
package app

import (
	"github.com/AlexLevus/telegram-bot/internal/repository"
	"github.com/AlexLevus/telegram-bot/internal/telegram_bot"
)

func Run() error {
	counterRepository, err := repository.NewRepository()
	telegram_bot.Start(counterRepository)

	return err
}
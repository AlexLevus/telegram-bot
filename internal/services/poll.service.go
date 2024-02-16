package services

import "github.com/AlexLevus/telegram-bot/internal/models"

type PollService interface {
	SavePoll(*models.SavePollRequest) error
}
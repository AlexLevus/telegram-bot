package services

import "github.com/AlexLevus/telegram-bot/internal/models"

type PollService interface {
	SavePoll(*models.SavePollRequest) error
	FindPollById(pollID string) (*models.DBPoll, error)
	UpdatePollStatus(pollID string, newStatus string) error
}
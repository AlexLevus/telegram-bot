package services

import "github.com/AlexLevus/telegram-bot/internal/models"

type SuggestionService interface {
	SaveSuggestion(*models.SaveSuggestionRequest) error
	FindSuggestionsByChatId(chatID int64) ([]*models.DBSuggestion, error)
	IsSuggestionExists(chatID int64, filmID int) bool
}
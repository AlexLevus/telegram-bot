package services

import "github.com/AlexLevus/telegram-bot/internal/models"

type ChatService interface {
	SaveChat(*models.SaveChatRequest) error
	UpdateChatMembersCount(chatID int64, membersCount int) error
	RemoveChatMember(chatID int64) error
	FindChatById(chatID int64) (*models.DBChat, error)
}
package services

import (
	"context"
	"time"

	"github.com/AlexLevus/telegram-bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SuggestionServiceImpl struct {
	suggestionCollection *mongo.Collection
	ctx            context.Context
}

func NewSuggestionService(suggestionCollection *mongo.Collection, ctx context.Context) SuggestionService {
	return &SuggestionServiceImpl{suggestionCollection, ctx}
}

func (c *SuggestionServiceImpl) SaveSuggestion(suggestion *models.SaveSuggestionRequest) error {
	suggestion.AddedAt = time.Now()

	_, err := c.suggestionCollection.InsertOne(c.ctx, suggestion)

	if err != nil {
		return err
	}

	return nil
}

func (c *SuggestionServiceImpl) IsSuggestionExists(chatID int64, filmID int) bool {
	query := bson.D{{Key: "chat_id", Value: chatID }, {Key: "film_id", Value: filmID }}

	var suggestion *models.DBSuggestion

	if err := c.suggestionCollection.FindOne(c.ctx, query).Decode(&suggestion); err != nil {
		return false
	}

	return true
}


func (c *SuggestionServiceImpl) FindSuggestionsByChatId(chatID int64) ([]*models.DBSuggestion, error) {
	query := bson.D{{Key: "chat_id", Value: chatID }}

	cursor, err := c.suggestionCollection.Find(c.ctx, query)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(c.ctx)

	var suggestions []*models.DBSuggestion

	for cursor.Next(c.ctx) {
		suggestion := &models.DBSuggestion{}
		err := cursor.Decode(suggestion)

		if err != nil {
			return nil, err
		}

		suggestions = append(suggestions, suggestion)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(suggestions) == 0 {
		return []*models.DBSuggestion{}, nil
	}

	return suggestions, nil

}

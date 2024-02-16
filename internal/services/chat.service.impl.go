package services

import (
	"context"
	"errors"
	"time"
	
	"github.com/AlexLevus/telegram-bot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatServiceImpl struct {
	chatCollection *mongo.Collection
	ctx            context.Context
}

func NewChatService(postCollection *mongo.Collection, ctx context.Context) ChatService {
	return &ChatServiceImpl{postCollection, ctx}
}

func (c *ChatServiceImpl) SaveChat(chat *models.SaveChatRequest) error {
	chat.AddedAt = time.Now()

	_, err := c.chatCollection.InsertOne(c.ctx, chat)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return c.UpdateChatMembersCount(chat.ChatId, chat.MembersCount)
		}
		return err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"chat_id": 1}, Options: opt}

	if _, err := c.chatCollection.Indexes().CreateOne(c.ctx, index); err != nil {
		return errors.New("could not create index for ID")
	}

	return nil
}

func (c *ChatServiceImpl) FindChatById(id int64) (*models.DBChat, error) {
	query := bson.D{{Key: "chat_id", Value: id }}

	var chat *models.DBChat

	if err := c.chatCollection.FindOne(c.ctx, query).Decode(&chat); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return chat, nil
}

func (c *ChatServiceImpl) UpdateChatMembersCount(chatID int64, membersCount int) error {
	query := bson.D{{Key: "chat_id", Value: chatID}}
	update := bson.M{"$set": bson.M{"members_count": membersCount }}

	_, err := c.chatCollection.UpdateOne(
		c.ctx,
		query,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *ChatServiceImpl) RemoveChatMember(chatID int64) error {
	chat, err := c.FindChatById(chatID)

	if err != nil {
		return err
	}

	query := bson.D{{Key: "chat_id", Value: chatID}}
	update := bson.M{"$set": bson.M{"members_count": chat.MembersCount - 1 }}

	_, err = c.chatCollection.UpdateOne(
		c.ctx,
		query,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}
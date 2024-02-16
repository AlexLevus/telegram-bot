package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddChatRequest struct {
	ChatId    int64    `json:"chat_id" bson:"chat_id" binding:"required"`
	MembersCount    int    `json:"members_count" bson:"members_count" binding:"required"`
	AddedAt          time.Time `json:"added_at" bson:"added_at"`
}

type DBChat struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId    	int64    `json:"chat_id" bson:"chat_id"`
	MembersCount    int    `json:"members_count" bson:"members_count"`
	AddedAt          time.Time          `json:"added_at" bson:"added_at"`
}
package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavePollRequest struct {
	ChatId    int64    `json:"chat_id" bson:"chat_id" binding:"required"`
	PollId    	int64    `json:"poll_id" bson:"poll_id"`
	AddedAt          time.Time `json:"added_at" bson:"added_at"`
}

type DBPoll struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PollId    	int64    `json:"poll_id" bson:"poll_id"`
	ChatId    	int64    `json:"chat_id" bson:"chat_id"`
	AddedAt          time.Time          `json:"added_at" bson:"added_at"`
}
package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaveSuggestionRequest struct {
	ChatId    int64    `json:"chat_id" bson:"chat_id" binding:"required"`
	FilmId    	int    `json:"film_id" bson:"film_id" binding:"required"`
	FilmName    	string    `json:"film_name" bson:"film_name" binding:"required"`
	FilmUrl    	string    `json:"film_url" bson:"film_url" binding:"required"`
	FilmRating    	float32    `json:"film_rating" bson:"film_rating" binding:"required"`
	FilmYear   	int    `json:"film_year" bson:"film_year" binding:"required"`
	AddedAt          time.Time `json:"added_at" bson:"added_at" binding:"required"`
}

type DBSuggestion struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId    int64    `json:"chat_id" bson:"chat_id"`
	FilmId    	int    `json:"film_id" bson:"film_id"`
	FilmName    	string    `json:"film_name" bson:"film_name"`
	FilmUrl    	string    `json:"film_url" bson:"film_url"`
	FilmRating    	float32    `json:"film_rating" bson:"film_rating"`
	FilmYear   	int    `json:"film_year" bson:"film_year"`
	AddedAt          time.Time `json:"added_at" bson:"added_at" binding:"required"`
}
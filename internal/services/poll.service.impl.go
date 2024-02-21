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

type PollServiceImpl struct {
	pollCollection *mongo.Collection
	ctx            context.Context
}

func NewPollService(postCollection *mongo.Collection, ctx context.Context) PollService {
	return &PollServiceImpl{postCollection, ctx}
}

func (c *PollServiceImpl) SavePoll(poll *models.SavePollRequest) error {
	poll.AddedAt = time.Now()

	_, err := c.pollCollection.InsertOne(c.ctx, poll)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return errors.New("poll with that title already exists")
		}
		return err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"poll_id": 1}, Options: opt}

	if _, err := c.pollCollection.Indexes().CreateOne(c.ctx, index); err != nil {
		return errors.New("could not create index for ID")
	}

	return nil
}

func (c *PollServiceImpl) FindPollById(id string) (*models.DBPoll, error) {
	query := bson.D{{Key: "poll_id", Value: id }}

	var poll *models.DBPoll

	if err := c.pollCollection.FindOne(c.ctx, query).Decode(&poll); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return poll, nil
}

func (c *PollServiceImpl) UpdatePollStatus(pollId string, newStatus string) error {
	query := bson.D{{Key: "poll_id", Value: pollId}}
	update := bson.M{"$set": bson.M{"status": newStatus }}

	_, err := c.pollCollection.UpdateOne(
		c.ctx,
		query,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

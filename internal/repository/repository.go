package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Repository struct {
	ctx        context.Context
	collection *mongo.Collection
}

type Counter struct {
	Value     int                `bson:"value,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updatedAt,omitempty"`
}

func NewRepository() (*Repository, error) {
	mongoDbUri, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		log.Fatal("Добавьте в файл .env uri к базе MongoDB")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	collection := client.Database("CounterDB").Collection("Counter")

	return &Repository{collection: collection, ctx: ctx}, nil
}

func (repository *Repository) GetCounter() (Counter, error) {
	cur, err := repository.collection.Find(repository.ctx, bson.D{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatalf("Error when close DB cursor")
		}
	}(cur, repository.ctx)

	var counters []Counter

	err = cur.All(repository.ctx, &counters)
	if err != nil {
		log.Fatalf("Error when get Counters from DB")
	}

	return counters[0], nil
}

func (repository *Repository) UpdateCounter(counter Counter) error {
	id, _ := primitive.ObjectIDFromHex("63692f15b50ce6ea336f9139")
	filter := bson.D{{"_id", id}}

	updatedAt := time.Now()

	update := bson.M{
		"$set": bson.M{
			"value":     counter.Value + 1,
			"updatedAt": updatedAt,
		},
	}

	_, err := repository.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	return nil
}

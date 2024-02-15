package repository

import (
	"context"
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

type Chat struct {
	Value     int                `bson:"value,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updatedAt,omitempty"`
}

func NewRepository() (*Repository, error) {
	mongoDbUri, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		log.Fatal("Добавьте в файл .env uri к базе MongoDB")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatal("Добавьте в файл .env название базы данных")
	}

	clientOptions := options.Client().ApplyURI(mongoDbUri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(dbName)
	chatCollection := database.Collection("chats")

	return &Repository{collection: chatCollection, ctx: ctx}, nil
}

func (repository *Repository) GetChat() (Chat, error) {
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

	var chats []Chat

	err = cur.All(repository.ctx, &chats)
	if err != nil {
		log.Fatalf("Error when get chats from DB")
	}

	return chats[0], nil
}

func (repository *Repository) UpdateChat(chat Chat) error {
	id, _ := primitive.ObjectIDFromHex("63692f15b50ce6ea336f9139")
	filter := bson.D{{"_id", id}}

	updatedAt := time.Now()

	update := bson.M{
		"$set": bson.M{
			"value":     chat.Value + 1,
			"updatedAt": updatedAt,
		},
	}

	_, err := repository.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	return nil
}

package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
}

func (*Repo) InsertMany(db *mongo.Client, data string) {
	collection := db.Database("manabie-webhook-test").Collection("email-event")
	docs := []interface{}{
		bson.D{{"data", data}},
	}
	ctx := context.Background()
	_, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
}

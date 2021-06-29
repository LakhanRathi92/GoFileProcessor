package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func StorePersonData(client *mongo.Client, person *Person) (result *mongo.InsertOneResult, err error) {
	collection := client.Database("MyApp").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err = collection.InsertOne(ctx, person)
	return result, err
}

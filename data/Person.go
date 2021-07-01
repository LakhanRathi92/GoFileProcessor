package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

//create a single person.
func (p *Person) Store(client *mongo.Client) (result *mongo.InsertOneResult, err error) {
	collection := client.Database("MyApp").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err = collection.InsertOne(ctx, p)
	return result, err
}

//reads all matching data where person's first name is given.
func Read(client *mongo.Client, firstname string) (personsFiltered []bson.M, err error) {
	personCollection := client.Database("MyApp").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filterCursor, err := personCollection.Find(ctx, bson.M{"firstname": firstname})

	if err != nil {
		return nil, err
	}

	err = filterCursor.All(ctx, &personsFiltered)
	if err != nil {
		return nil, err
	}
	return personsFiltered, err
}

//update a single person.

//delete a single person.

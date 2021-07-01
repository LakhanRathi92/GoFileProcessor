package data

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Store a new random created user, read the user, and then finally delete it.

var client *mongo.Client
var firstname string
var lastname string

var _id primitive.ObjectID

//initializes database, and user data.

func setup() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatalf("init test failed... %s", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //time out for connection attempts.
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("init test failed... %s", err)
	}

	len := 5
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("init test failed... %s", err)
	}

	firstname = fmt.Sprintf("%X", b)
	lastname = fmt.Sprintf("%X", b)
}

//creates the person and stores it
func TestStore(t *testing.T) {
	setup()
	p := Person{
		Firstname: firstname,
		Lastname:  lastname,
	}
	result, err := p.Store(client)
	_id = result.InsertedID.(primitive.ObjectID)

	if err != nil {
		t.Errorf("Failed store test %s", err)
	}

	if p.Firstname != firstname {
		t.Errorf("Failed store test %s", err)
	}

	t.Log("Person was successfully stored : "+p.Firstname+" "+p.Lastname+" inserted id: ", _id)

	if err != nil {
		t.Errorf("Failed store test %s", err)
	}

}

//reads the newly created person.

func TestRead(t *testing.T) {
	personsFiltered, err := Read(client, firstname)
	if err != nil {
		t.Errorf("Failed read test :%s", err)
	}

	for _, value := range *personsFiltered {
		if value.ID == _id {
			t.Log("person exists :", value.ID)
			break
		}
	}
}

func TestTearDown(t *testing.T) {
	collection := client.Database("MyApp").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": _id})

	if err != nil {
		t.Log(err)
	}
	t.Log("Person was successfully removed : ", _id)
}

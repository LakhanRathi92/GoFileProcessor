package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//simple handler
type QueryHandler struct {
	l      *log.Logger
	client *mongo.Client
}

//new file handler creation function
func NewQueryHandler(l *log.Logger, client *mongo.Client) *QueryHandler {
	return &QueryHandler{l, client}
}

func (h *QueryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Add("content-type", "application/json")

		personCollection := h.client.Database("MyApp").Collection("person")
		persons, ok := r.URL.Query()["person"]

		if !ok || len(persons[0]) < 1 {
			log.Println("Url Param 'person' is missing")
			return
		}

		person := persons[0]
		log.Println("Url Param 'person' request: " + person)

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		filterCursor, err := personCollection.Find(ctx, bson.M{"firstname": person})

		if err != nil {
			log.Fatal(err)
		}

		var personsFiltered []bson.M
		if err = filterCursor.All(ctx, &personsFiltered); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(rw).Encode(personsFiltered)

	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

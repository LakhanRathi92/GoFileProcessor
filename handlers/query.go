package handlers

//handle incoming GET request to retrieve data.

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LakhanRathi92/GoFileProcessor/data"
	"go.mongodb.org/mongo-driver/mongo"
)

//simple handler
type QueryHandler struct {
	l      *log.Logger
	client *mongo.Client
}

//new Query handler creation function
func NewQueryHandler(l *log.Logger, client *mongo.Client) *QueryHandler {
	return &QueryHandler{l, client}
}

func (h *QueryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Add("content-type", "application/json")
		persons, ok := r.URL.Query()["person"]

		//check query parameter.
		if !ok || len(persons[0]) < 1 {
			h.l.Println("Url Param 'person' is missing")
			json.NewEncoder(rw).Encode("Url Param 'person' is missing")
			return
		}

		person := persons[0]
		h.l.Println("Url Param 'person' request: " + person)

		personsFiltered, err := data.Read(h.client, person)
		if err != nil {
			h.l.Print("no users found. ")
			json.NewEncoder(rw).Encode("no users found")
		}
		if personsFiltered != nil {
			json.NewEncoder(rw).Encode(personsFiltered)
		} else {
			json.NewEncoder(rw).Encode("no users found with matching first name: " + person)
		}
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

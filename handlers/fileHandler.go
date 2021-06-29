package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LakhanRathi92/GoFileProcessor/data"
	"go.mongodb.org/mongo-driver/mongo"
)

//simple handler
type FileHandler struct {
	l      *log.Logger
	client *mongo.Client
}

//new file handler creation function
func NewFileHandler(l *log.Logger, client *mongo.Client) *FileHandler {
	return &FileHandler{l, client}
}

//interface impl
func (h *FileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fixCrossOrigin(rw, r)

		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()

		h.l.Printf("Uploaded File: %+v\n", header.Filename)
		h.l.Printf("File Size: %+v\n", header.Size)
		h.l.Printf("MIME Header: %+v\n", header.Header)

		var person data.Person
		json.NewDecoder(file).Decode(&person)
		result, err := data.StorePersonData(h.client, &person)

		json.NewEncoder(rw).Encode(result)

	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func fixCrossOrigin(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("content-type", "application/json")
	rw.Header().Set("Access-Control-Allow-Headers:", "*")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")

	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}
}

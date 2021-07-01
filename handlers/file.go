package handlers

//reads data from file and creates person from it and calls store to store it in database.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LakhanRathi92/GoFileProcessor/data"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		rw.Header().Add("content-type", "application/json")

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

		var p data.Person
		json.NewDecoder(file).Decode(&p)
		result, err := p.Store(h.client)

		if err != nil {
			h.l.Fatal(err)
		}

		json.NewEncoder(rw).Encode(result)

	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func fixCrossOrigin(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Set("Access-Control-Allow-Headers:", "*")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")

	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}
}

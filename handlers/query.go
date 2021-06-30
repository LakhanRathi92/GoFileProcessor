package handlers

import (
	"log"
	"net/http"

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
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

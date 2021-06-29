package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LakhanRathi92/GoFileProcessor/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	l := log.New(os.Stdout, "file-processor", log.LstdFlags)
	l.Println("File processor started...")

	initDbConnection(l)
	startServer(l)
	close(l)
}

func close(l *log.Logger) {
	if client == nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Disconnect(ctx)
	if err != nil {
		l.Fatal(err)
	}
}

func initDbConnection(l *log.Logger) {
	l.Println("Initializing db connection...")

	var err error

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //time out for connection attempts.

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func startServer(l *log.Logger) {
	l.Println("Initializing server...")

	var s *http.Server

	sm := http.NewServeMux()

	fileUploadHandler := handlers.NewFileHandler(l, client)
	sm.Handle("/upload/file/", fileUploadHandler)

	sm.Handle("/", http.FileServer(http.Dir("./html")))

	s = &http.Server{
		Addr:         "127.0.0.1:9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//register OS signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Recieved terminate signal, graceful shutdown.... ", sig)

	//graceful shutdown. Wait for 30 seconds on existing handlers.
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}

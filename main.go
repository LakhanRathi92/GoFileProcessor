package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LakhanRathi92/GoFileProcessor/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//mongo.Client type is safe for concurrent use: https://developer.mongodb.com/community/forums/t/best-way-to-refactor-connection-overhead-from-my-handler-functions/3672
var client *mongo.Client

func main() {
	l := log.New(os.Stdout, "file-processor", log.LstdFlags)
	l.Println("File processor started...")

	serverInterruptSig := make(chan os.Signal, 1)
	signal.Notify(serverInterruptSig, syscall.SIGINT, syscall.SIGTERM) //Control-C interrupt,  sigterm sent by windows/process.

	//initialize DB
	err := initDbConnection(l)
	if err != nil {
		l.Fatal(err) //Program exits immediately.
	}

	//non-blocking start server.
	go startServer(l)

	//wait for the server to be stopped.
	sig := <-serverInterruptSig
	l.Printf("Server stopped with: %v \n", sig)

	disconnectDb(l)
}

func startServer(l *log.Logger) {
	l.Println("Initializing server...")
	fileUploadHandler := handlers.NewFileHandler(l, client)
	queryHandler := handlers.NewQueryHandler(l, client)

	sm := http.NewServeMux()

	//handler for static files.
	sm.Handle("/", http.FileServer(http.Dir("./html")))

	sm.Handle("/create/", fileUploadHandler)
	sm.Handle("/query", queryHandler)

	//todo delete and update handlers.

	http.ListenAndServe("127.0.0.1:9090", sm)
}

func initDbConnection(l *log.Logger) error {
	l.Println("Initializing db connection...")

	var err error

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //time out for connection attempts.

	err = client.Connect(ctx)

	return err
}

func disconnectDb(l *log.Logger) {
	l.Println("Closing db connection... ")
	if client == nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Disconnect(ctx)
	if err != nil {
		l.Fatal(err)
	}
	l.Println("Db connection closed...")
}

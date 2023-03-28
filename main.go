package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo(ctx context.Context) (context.Context, *mongo.Client) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://anhhuu:ujDLYrIboxq9jCpR@cluster0.16sjn.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return ctx, client
}

func insert(client *mongo.Client, data string) {
	collection := client.Database("manabie-webhook-test").Collection("email-event")
	docs := []interface{}{
		bson.D{{"data", data}},
	}
	ctx := context.Background()
	res, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)
}

var clientDB *mongo.Client

func main() {
	http.HandleFunc("/", handlerMain)

	http.HandleFunc("/url", handler)

	port := os.Getenv("PORT")

	ctx := context.Background()

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	_, clientDB = connectMongo(ctx)

	fmt.Printf("Starting server at port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handlerMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Fprintf(w, "body:\n%s\n", body)
	log.Printf("body:\n%s\n", body)

	insert(clientDB, string(body))
}

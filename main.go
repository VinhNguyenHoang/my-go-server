package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handlerMain)

	http.HandleFunc("/url", handler)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

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
}

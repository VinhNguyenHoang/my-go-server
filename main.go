package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlerMain)

	http.HandleFunc("/url", handler)
	fmt.Print("Starting server at port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
}

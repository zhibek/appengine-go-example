package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Save new "Demo" entity to datastore
	newDemo := &Demo{Name: "Test"}
	newDemo, err := CreateDemo(ctx, newDemo)
	if err != nil {
		log.Printf("Error creating entity: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Created entity with ID %s", newDemo.ID)

	// Fetch "Demo" entity from datastore
	retrievedDemo, err := FetchDemo(ctx, newDemo.ID)
	if err != nil {
		log.Printf("Error fetching entity: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched entity with ID %s", newDemo.ID)

	fmt.Fprint(w, retrievedDemo.Name)
}

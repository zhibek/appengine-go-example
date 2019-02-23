package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/zhibek/appengine-go-example/example"
)

func main() {
	r := chi.NewRouter()

	r.Route("/example", func(r chi.Router) {
		r.Get("/", example.ListHandler)
		r.Post("/", example.CreateHandler)
		r.Route("/{exampleId}", func(r chi.Router) {
			r.Use(example.RouteIdParser)
			r.Get("/", example.GetHandler)
			r.Put("/", example.UpdateHandler)
			r.Delete("/", example.DeleteHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

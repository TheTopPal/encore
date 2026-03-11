package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

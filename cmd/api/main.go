package main

import (
	"log"
	"net/http"
	"os"

	"csv-importer/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Post("/upload", handlers.UploadHandler)

	port := os.Getenv("PORT")

	log.Printf("Starting on PORT=%s", port)
	http.ListenAndServe(":"+port, router)
}

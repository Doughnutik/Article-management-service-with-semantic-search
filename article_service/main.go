package main

import (
	"log"
	"net/http"

	"article_service/api"
	handler "article_service/internal/handler"
)

func main() {
	// Create service instance.
	handler := handler.NewArticleHandler()
	// Create generated server.
	srv, err := api.NewServer(handler)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}

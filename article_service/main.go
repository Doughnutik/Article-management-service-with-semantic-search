package main

import (
	"log"
	"net/http"

	"article_service/api"
	"article_service/internal/grpc_client"
	handler "article_service/internal/handler"
)

func main() {
	// Создаём gRPC клиент для Python SearchService
	searchClient, err := grpc_client.NewSearchClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer searchClient.Close()

	// Создаём REST handler с gRPC клиентом
	h := handler.NewArticleHandler(searchClient)

	// Создаём сервер, сгенерированный ogen
	srv, err := api.NewServer(h)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем REST API
	log.Println("Starting REST server on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}

// curl -i -X 'POST' \
//   'http://localhost:8080/articles/search' \
//   -H 'accept: application/json' \
//   -H 'Content-Type: application/json' \
//   -d '{
//   "query": "string",
//   "limit": 3
// }'
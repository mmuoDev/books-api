package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"books-api/internal/app"

	"github.com/mmuoDev/commons/mongo"
)

func main() {

	port := os.Getenv("MONGO_PORT")
	provideMongoDB, err := mongo.NewConfigFromEnvVars().ToProvider(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(provideMongoDB)
	log.Println(fmt.Sprintf("Starting server on port:%s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), a.Handler()))
}

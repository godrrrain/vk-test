package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	"vktest/src/handler"
	"vktest/src/storage"
	"vktest/src/tools"
)

func main() {
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		"postgres", 5432, "program", "movies", "test")
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	handler := handler.NewHandler(psqlDB)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/get/movies", tools.RequestLogger(handler.GetMovies))
	mux.HandleFunc("/api/v1/get/actors", tools.RequestLogger(handler.GetActors))

	mux.HandleFunc("/api/v1/post/movies", tools.RequestAuth(handler.CreateMovie))
	mux.HandleFunc("/api/v1/post/actors", tools.RequestAuth(handler.CreateActor))

	mux.HandleFunc("/api/v1/delete/movies", tools.RequestAuth(handler.DeleteMovie))
	mux.HandleFunc("/api/v1/delete/actors", tools.RequestAuth(handler.DeleteActor))

	mux.HandleFunc("/api/v1/upd/actors", tools.RequestAuth(handler.UpdateActor))
	mux.HandleFunc("/api/v1/upd/movie", tools.RequestAuth(handler.UpdateMovie))
	mux.HandleFunc("/api/v1/search/movies", tools.RequestLogger(handler.SearchMovies))

	corsCustom := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	corsHandler := corsCustom.Handler(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatal(err)
	}
}

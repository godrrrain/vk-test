package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/api/v1/get/movies", tools.RequestLogger(handler.GetMovies))
	http.HandleFunc("/api/v1/get/actors", tools.RequestLogger(handler.GetActors))

	http.HandleFunc("/api/v1/post/movies", tools.RequestAuth(handler.CreateMovie))
	http.HandleFunc("/api/v1/post/actors", tools.RequestAuth(handler.CreateActor))

	http.HandleFunc("/api/v1/delete/movies", tools.RequestAuth(handler.DeleteMovie))
	http.HandleFunc("/api/v1/delete/actors", tools.RequestAuth(handler.DeleteActor))

	http.HandleFunc("/api/v1/upd/actors", tools.RequestAuth(handler.UpdateActor))
	http.HandleFunc("/api/v1/upd/movie", tools.RequestAuth(handler.UpdateMovie))
	http.HandleFunc("/api/v1/search/movies", tools.RequestLogger(handler.SearchMovies))

	port := "8080"
	// Запуск сервера
	log.Println("Starting server on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

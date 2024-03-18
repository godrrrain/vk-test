package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"vktest/src/storage"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type MovieResponse struct {
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Release_date string              `json:"release_date"`
	Rating       int                 `json:"rating"`
	Actors       []storage.ActorName `json:"actors"`
}

type CreateMovieRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Release_date string `json:"release_date"`
	Rating       int    `json:"rating"`
	Actors       []int  `json:"actors"`
}

type CreateActorRequest struct {
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UpdateActorRequest struct {
	ID       int    `json:"id" binding:"required"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UpdateMovieRequest struct {
	ID           int    `json:"id" binding:"required"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Release_date string `json:"release_date"`
	Rating       int    `json:"rating"`
}

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {

	sortField := r.URL.Query().Get("sort")
	if sortField == "" || sortField == "rating" {
		sortField = "rating DESC"
	}
	if sortField == "title" || sortField == "name" {
		sortField = "title ASC"
	}
	if sortField == "release_date" || sortField == "release" || sortField == "date" {
		sortField = "release_date DESC"
	}

	movies, err := h.storage.GetMovies(context.Background(), sortField)

	if err != nil {
		log.Printf("failed to get movies %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (h *Handler) GetActors(w http.ResponseWriter, r *http.Request) {

	actors, err := h.storage.GetActors(context.Background())

	if err != nil {
		log.Printf("failed to get actors %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}

func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "error: Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var actorBody CreateActorRequest

	if err := json.NewDecoder(r.Body).Decode(&actorBody); err != nil {
		http.Error(w, "error: Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.storage.CreateActor(context.Background(), actorBody.Name, actorBody.Gender, actorBody.Birthday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully created",
	})
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "error: Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var movieBody CreateMovieRequest

	if err := json.NewDecoder(r.Body).Decode(&movieBody); err != nil {
		http.Error(w, "error: Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.storage.CreateMovie(context.Background(), movieBody.Title, movieBody.Description, movieBody.Release_date, movieBody.Rating, movieBody.Actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully created",
	})
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "error: Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	actorIdStr := r.URL.Query().Get("id")
	if actorIdStr == "" {
		http.Error(w, "error: Actor ID is required", http.StatusBadRequest)
		return
	}

	actorId, err := strconv.Atoi(actorIdStr)
	if err != nil {
		http.Error(w, "error: Invalid Actor ID", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteActor(context.Background(), actorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully deleted",
	})
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "error: Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	movieIdStr := r.URL.Query().Get("id")
	if movieIdStr == "" {
		http.Error(w, "error: Actor ID is required", http.StatusBadRequest)
		return
	}

	movieId, err := strconv.Atoi(movieIdStr)
	if err != nil {
		http.Error(w, "error: Invalid Actor ID", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteMovie(context.Background(), movieId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully deleted",
	})
}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "error: Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var actorBody UpdateActorRequest

	if err := json.NewDecoder(r.Body).Decode(&actorBody); err != nil {
		http.Error(w, "error: Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.storage.UpdateActor(context.Background(), actorBody.ID, actorBody.Name, actorBody.Gender, actorBody.Birthday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully updated",
	})
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "error: Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var movieBody UpdateMovieRequest

	if err := json.NewDecoder(r.Body).Decode(&movieBody); err != nil {
		http.Error(w, "error: Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.storage.UpdateMovie(context.Background(), movieBody.ID, movieBody.Title, movieBody.Description, movieBody.Release_date, movieBody.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "successfully updated",
	})
}

func (h *Handler) SearchMovies(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	sortField := r.URL.Query().Get("sort")
	if sortField == "" || sortField == "rating" {
		sortField = "rating DESC"
	}
	if sortField == "title" || sortField == "name" {
		sortField = "title ASC"
	}
	if sortField == "release_date" || sortField == "release" || sortField == "date" {
		sortField = "release_date DESC"
	}

	movies, err := h.storage.GetMovies(context.Background(), sortField)
	if err != nil {
		log.Printf("failed to get movies %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var acceptMovies []storage.MovieInfo

	for _, v := range movies {
		checkActorName := false
		for _, val := range v.Actors {
			if strings.Contains(val.Name, search) {
				checkActorName = true
				break
			}
		}

		if strings.Contains(v.Title, search) || checkActorName {
			acceptMovies = append(acceptMovies, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acceptMovies)
}

// func MovieToResponse(movie storage.MovieInfo) MovieResponse {
// 	return MovieResponse{
// 		Title:        movie.Title,
// 		Description:  movie.Description,
// 		Release_date: movie.Release_date,
// 		Rating:       movie.Rating,
// 		Actors:       movie.Actors,
// 	}
// }

// func MoviesToResponse(movies []storage.MovieInfo) []MovieResponse {
// 	if movies == nil {
// 		return nil
// 	}

// 	res := make([]MovieResponse, len(movies))

// 	for index, value := range movies {
// 		res[index] = MovieToResponse(value)
// 	}

// 	return res
// }

package router

import (
	"github.com/gorilla/mux"
	"github.com/Madhav-M01/mangodb/controller"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/movies", controller.GetMyMovies).Methods("GET")
	r.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("POST")
	r.HandleFunc("/api/movie/{id}", controller.DeleteOneMovieHandler).Methods("DELETE") // http://localhost:4000/api/movie/68541c82b7af20636459ec42
	r.HandleFunc("/api/movies", controller.DeleteAllMoviesHandler).Methods("DELETE")
	r.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST") 


	return r
}

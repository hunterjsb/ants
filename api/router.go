package api

import (
	"ants/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/start", handlers.Start).Methods("GET")

	r.HandleFunc("/colonies/{colonyIndex}/spawn/{antType}", handlers.Spawn).Methods("POST")

	amw := AuthenticationMiddleware{}
	r.Use(amw.Middleware)
	r.Use(loggingMiddleware)

	return r
}

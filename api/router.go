package api

import (
	"ants/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", handlers.AntHandler).Methods("GET")

	amw := AuthenticationMiddleware{}
	r.Use(amw.Middleware)

	return r
}

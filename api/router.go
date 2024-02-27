package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ants!\n"))
}

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler).Methods("GET")

	amw := authenticationMiddleware{tokenUsers: make(map[string]string)}
	amw.Populate()
	r.Use(amw.Middleware)

	return r
}

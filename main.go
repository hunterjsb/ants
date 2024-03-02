package main

import (
	"log"
	"net/http"

	"ants/api"
	"ants/world"
)

func main() {
	// Create a new world
	world.Init()

	// Set up the web server
	r := api.RegisterRoutes()

	// Start listening on port 8000
	log.Fatal(http.ListenAndServe(":8000", r))
}

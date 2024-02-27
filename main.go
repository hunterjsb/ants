package main

import (
	"log"
	"net/http"

	"ants/api"
	"ants/world"
)

func main() {
	world.Init()
	// queen := ant.Ant{Tile: w.Tiles[1][1], Type: ant.Queen, MoveSpeed: 0}

	r := api.RegisterRoutes()
	log.Fatal(http.ListenAndServe(":8000", r))
}

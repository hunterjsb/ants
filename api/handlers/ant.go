package handlers

import (
	"ants/ant"
	"ants/world"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func Start(w http.ResponseWriter, r *http.Request) {
	randomX := rand.Intn(world.Width)
	randomY := rand.Intn(world.Height)
	user, _ := r.Context().Value("user").(ant.User)

	colony := ant.NewColony(&user, world.OverWorld.Tiles[randomY][randomX])
	fmt.Println("New Colony created for", colony.Owner.Name)

	// Marshal the newQueen to JSON
	jsonResp, err := json.Marshal(*(colony.Queen))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

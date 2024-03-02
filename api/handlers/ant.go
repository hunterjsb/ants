package handlers

// import (
// 	"ants/ant"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func Spawn(w http.ResponseWriter, r *http.Request) {
// 	user, ok := ant.UserFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
// 		return
// 	}

// 	params := mux.Vars(r)
// 	queen := user.Colonies[params["colonyIndex"]]

// 	fmt.Println("New Colony created for", colony.Owner.Name)

// 	// Marshal the newQueen to JSON
// 	jsonResp, err := json.Marshal(*(colony.Queen))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the Content-Type and write the JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonResp)
// }

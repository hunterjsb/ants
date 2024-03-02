package handlers

import (
	"ants/ant"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Spawn(w http.ResponseWriter, r *http.Request) {
	user, ok := ant.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
		return
	}

	fmt.Println("USER ", user)
	params := mux.Vars(r)
	colIdxStr := params["colonyIndex"]
	colonyIndex, err := strconv.Atoi(colIdxStr)
	if err != nil || colonyIndex < 0 || colonyIndex > (len(user.Colonies)-1) {
		http.Error(w, fmt.Sprintf("Invalid colony index %s", colIdxStr), http.StatusBadRequest)
		return
	}

	colony := user.Colonies[colonyIndex]
	fmt.Println(colony)

	// fmt.Println("New Colony created for", colony.Owner.Name)

	// // Marshal the newQueen to JSON
	// jsonResp, err := json.Marshal(*(colony.Queen))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Set the Content-Type and write the JSON response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(jsonResp)
}

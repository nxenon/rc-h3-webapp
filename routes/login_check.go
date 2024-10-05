package routes

import (
	"encoding/json"
	"net/http"
)

// /api/login_check -> checks if user is logged in
// if user is logged in returns {'logged_in': true}, otherwise {'logged_in': false}

func LoginCheckRouteHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]bool{
		"logged_in": false,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

}

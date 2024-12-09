package routes

import (
	"encoding/json"
	"net/http"
	"waitlist-golang/models"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	
	// get the ip address of the client
	ip := r.RemoteAddr

	if r.Method == "GET" {
		responseModel := models.Ping{
			Response: 200,
			Message: "GET request received",
			IP: ip,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(responseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
		
	if r.Method == "POST" {
		responseModel := models.Ping{
			Response: 200,
			Message: "POST request received",
			IP: ip,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(responseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} 

	if r.Method != "GET" && r.Method != "POST" {
		responseModel := models.Ping{
			Response: 405,
			Message: "Method not allowed",
			IP: ip,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(responseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}



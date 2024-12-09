package routes

import (
	"fmt"
	"net/http"
	"waitlist-golang/database"
)

func PingDB(w http.ResponseWriter, r *http.Request) {
	err := database.DBPool.Ping(r.Context())

	if err != nil {
		http.Error(w, "Database not reachable", http.StatusInternalServerError)
		return
	}
	// set the response code
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Database is reachable")
}


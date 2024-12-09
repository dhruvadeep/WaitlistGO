package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"waitlist-golang/database"
	"waitlist-golang/utils"
)

type InfoUser struct {
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	IPAddress    string `json:"ipaddress"`
	RefBy        string `json:"ref_by"`
}

type InfoFromClient struct {
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	RefBy        string `json:"ref_by"`
}

type Response struct {
	Message string `json:"message"`
}


func SaveUserToDB(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var info InfoFromClient
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// check if all fields are present
	if info.Name == "" || info.EmailAddress == "" || info.RefBy == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	// save the user to the database

	var user InfoUser
	user.Name = info.Name
	user.EmailAddress = info.EmailAddress
	user.IPAddress = ip
	user.RefBy = info.RefBy

	// check if the email is correct
	// use a regex to check if the email is correct
	// if not return an error
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(user.EmailAddress) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Response{Message: "Invalid email address"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// check if email already exists using exists query
	var exists bool
	err = database.DBPool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, user.EmailAddress).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if email exists: %v", err)
		http.Error(w, fmt.Sprintf("Failed to check if email exists: %v", err), http.StatusInternalServerError)
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Response{Message: "Email already exists"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}




	// insert into the db
	_, err = database.DBPool.Exec(ctx, `
    INSERT INTO users (name, email, ipaddress, ref_by)
    VALUES ($1, $2, $3, $4)
	`, user.Name, user.EmailAddress, user.IPAddress, user.RefBy)

	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save user to database: %v", err), http.StatusInternalServerError)
		return
	}

	// insert into referrals table
	statement := `INSERT INTO referrals (ref_email, ref_to_email, ipaddress) VALUES ($1, $2, $3)`
	_, err = database.DBPool.Exec(ctx, statement, user.RefBy, user.EmailAddress, user.IPAddress)
	if err != nil {
		log.Printf("Error inserting referral into database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save referral to database: %v", err), http.StatusInternalServerError)
		return
	}

	// SEND EMAIL with go it becomes a async operation
	go utils.SendEmail("Dhruvadeep <waitlist@mail.dhruvadeep.dev>", user.EmailAddress, user.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{Message: "User saved successfully"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


}
package main

import (
	"fmt"
	"net/http"
	"waitlist-golang/database"
	"waitlist-golang/routes"
)

func main() {
	// CALL THE DB CALLS
	connURL := database.GetDatabaseURL()
	if connURL == "" {
		fmt.Println("Failed to get database URL")
		return
	}

	// initialize the database
	dbManager, err := database.NewDBManager(connURL)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	defer dbManager.Close()

	// REFRESH THE CONNECTION VERY 1HR
	dbManager.RefreshConnection(connURL)

	// DB MIGRATIONS
	database.RunMigrations(dbManager.GetPool())

	fmt.Println("Database initialized successfully")


	server := http.NewServeMux()

	server.HandleFunc("/ping", routes.Ping)
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	server.HandleFunc("/pingdb", routes.PingDB)
	server.HandleFunc("/signup", routes.SaveUserToDB)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe("localhost:8080", server)
}
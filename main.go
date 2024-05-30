package main

import (
	"fmt"
	"log"
	"net/http"

	"test1/db"
	"test1/handlers"
)

func main() {
	// Initialize the database
	db.InitDB()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Route setup
	http.HandleFunc("/protected", handlers.BasicAuth(handlers.ProtectedRouteHandler))
	http.HandleFunc("/create_person", handlers.CreatePersonHandler)
	http.HandleFunc("/update_person_age", handlers.UpdatePersonAgeHandler)

	// Start server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

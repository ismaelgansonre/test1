package handlers

import (
	"encoding/json"
	"net/http"

	"test1/db"
	"test1/models"
	"test1/utils"
)

// Handler for creating a new person
func CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var p models.Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db.Mutex.Lock() // Lock the mutex before accessing the database
	defer db.Mutex.Unlock() // Unlock the mutex after accessing the database

	// Check if the person already exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM person WHERE name = ?)", p.Name).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Person already exists", http.StatusConflict)
		return
	}

	_, err = db.DB.Exec("INSERT INTO person (name, age) VALUES (?, ?)", p.Name, p.Age)
	if err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, p)
}

// Handler for updating the age of the person
func UpdatePersonAgeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var update models.Person
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db.Mutex.Lock() // Lock the mutex before accessing the database
	defer db.Mutex.Unlock() // Unlock the mutex after accessing the database

	// Check if the person exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM person WHERE name = ?)", update.Name).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec("UPDATE person SET age = ? WHERE name = ?", update.Age, update.Name)
	if err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, update)
}

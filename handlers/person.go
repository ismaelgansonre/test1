package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.etcd.io/bbolt"
	"test1/db"
	"test1/models"
	"test1/utils"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
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

	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	err = db.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("persons"))
		if b.Get([]byte(p.Name)) != nil {
			return &utils.AppError{StatusCode: http.StatusConflict, Message: "Person already exists"}
		}
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(p.Name), data)
	})

	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			fmt.Printf("%sError: %s (status code: %d)%s\n", Red, appErr.Message, appErr.StatusCode, Reset)
			http.Error(w, appErr.Message, appErr.StatusCode)
		} else {
			fmt.Printf("%sError: Failed to create person (status code: %d)%s\n", Red, http.StatusInternalServerError, Reset)
			http.Error(w, "Failed to create person", http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("%sUser added successfully%s\n", Green, Reset)
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

	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	err = db.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("persons"))
		personData := b.Get([]byte(update.Name))
		if personData == nil {
			return &utils.AppError{StatusCode: http.StatusNotFound, Message: "Person not found"}
		}

		var person models.Person
		if err := json.Unmarshal(personData, &person); err != nil {
			return err
		}
		person.Age = update.Age

		data, err := json.Marshal(person)
		if err != nil {
			return err
		}
		return b.Put([]byte(person.Name), data)
	})

	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			fmt.Printf("%sError: %s (status code: %d)%s\n", Red, appErr.Message, appErr.StatusCode, Reset)
			http.Error(w, appErr.Message, appErr.StatusCode)
		} else {
			fmt.Printf("%sError: Failed to update person (status code: %d)%s\n", Red, http.StatusInternalServerError, Reset)
			http.Error(w, "Failed to update person", http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("%sUser updated successfully%s\n", Green, Reset)
	utils.SendJSONResponse(w, http.StatusOK, update)
}

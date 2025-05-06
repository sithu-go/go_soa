package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sithu-go/go-soa/common/models"
)

var users = []models.User{}

func init() {
	// Add some sample users
	users = append(users, models.User{
		ID:        "1",
		Username:  "john_doe",
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
	})
	users = append(users, models.User{
		ID:        "2",
		Username:  "jane_smith",
		Email:     "jane@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
	})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, user := range users {
		if user.ID == params["id"] {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	user.ID = uuid.New().String()
	users = append(users, user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func validateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, user := range users {
		if user.ID == params["id"] {
			json.NewEncoder(w).Encode(map[string]bool{"valid": true})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{"valid": false})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}/validate", validateUser).Methods("GET")

	log.Println("User Service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

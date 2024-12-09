

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

var db *sql.DB

type User struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Role     string `json:"role"`
}

type Chef struct {
	ChefID int    `json:"chef_id"`
	UserID int    `json:"user_id"`
	Bio    string `json:"bio"`
	Rating float64 `json:"rating"`
}

// Initialize the database connection
func initDB() {
	var err error
	connStr := "postgres://user:password@db:5432/userservice?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Users (
			user_id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(100) NOT NULL,
			location VARCHAR(100),
			role VARCHAR(10) DEFAULT 'USER',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS Chefs (
			chef_id SERIAL PRIMARY KEY,
			user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
			bio TEXT,
			rating DECIMAL(3, 2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS Follow (
			follow_id SERIAL PRIMARY KEY,
			user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
			chef_id INT REFERENCES Chefs(chef_id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS FavoriteDishes (
			favorite_id SERIAL PRIMARY KEY,
			user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
			dish_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

// Handlers

// GET: Retrieve a user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT user_id, email, name, location, role FROM Users WHERE user_id = $1", userID).
		Scan(&user.UserID, &user.Email, &user.Name, &user.Location, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GET: Retrieve all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT user_id, email, name, location, role FROM Users")
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Email, &user.Name, &user.Location, &user.Role)
		if err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST: Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set default role to "USER" if it's missing
	if user.Role == "" {
		user.Role = "USER"
	}

	// Insert the new user into the database
	err := db.QueryRow(
		"INSERT INTO Users(email, name, location, role) VALUES($1, $2, $3, $4) RETURNING user_id",
		user.Email, user.Name, user.Location, user.Role).Scan(&user.UserID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


// DELETE: Delete a user by ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	_, err = db.Exec("DELETE FROM Users WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Return 204 No Content
}



// PUT: Update a user (only update fields provided, do not change user_id)
// PUT: Update a user (only update fields provided, do not change user_id)
func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve current user data
	var existingUser User
	err = db.QueryRow("SELECT user_id, email, name, location, role FROM Users WHERE user_id = $1", userID).
		Scan(&existingUser.UserID, &existingUser.Email, &existingUser.Name, &existingUser.Location, &existingUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	// Parse incoming update fields
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare updated fields (only update provided fields)
	// If a field is empty, we keep the existing value.
	updatedLocation := existingUser.Location
	if user.Location != "" {
		updatedLocation = user.Location
	}

	updatedName := existingUser.Name
	if user.Name != "" {
		updatedName = user.Name
	}

	updatedRole := existingUser.Role
	if user.Role != "" {
		updatedRole = user.Role
	}

	// Update user data in the database
	_, err = db.Exec("UPDATE Users SET name = $1, location = $2, role = $3, updated_at = CURRENT_TIMESTAMP WHERE user_id = $4",
		updatedName, updatedLocation, updatedRole, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return the updated user data (only fields that were modified)
	existingUser.Name = updatedName
	existingUser.Location = updatedLocation
	existingUser.Role = updatedRole

	// Ensure the response includes all fields, with unmodified fields returned as they were
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingUser)
}


func main() {
	initDB()

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{user_id}", getUser).Methods("GET")
	router.HandleFunc("/users/{user_id}", updateUser).Methods("PUT")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{user_id}", deleteUser).Methods("DELETE")

	log.Println("User Service is running on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	UserID    int `json:"userID"`
	ExpiresAt time.Time
	jwt.StandardClaims
}
type auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type authResult struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	result := authResult{Success: false}
	w.Header().Set("Content-Type", "application/json")
	defer json.NewEncoder(w).Encode(&result)

	var creds auth
	json.NewDecoder(r.Body).Decode(&creds)

	// If no credentials, bail.
	if creds.Email == "" || creds.Password == "" {
		return
	}

	// Get user from DB by email.
	user := getUserByEmail(creds.Email)

	// If user exists, bail
	if user.ID != 0 {
		return
	}
	// Hash password and create new record
	hash, err := HashPassword(creds.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	insertUser(creds.Email, hash)
	user = getUserByEmail(creds.Email)
	// Create the Claims for our new user.
	claims := JWTClaims{
		user.ID,
		time.Now().Add(48 * time.Hour),
		jwt.StandardClaims{},
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(globalConfig.SigningSecret))
	if err != nil {
		log.Fatal(err)
	}
	result.Success = true
	result.Token = ss
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := authResult{Success: false}
	defer json.NewEncoder(w).Encode(&result)

	var creds auth
	json.NewDecoder(r.Body).Decode(&creds)
	user := getUserByEmail(creds.Email)
	if user.ID == 0 {
		return
	}
	if CheckPasswordHash(creds.Password, user.PasswordHash) {
		// Create the Claims
		claims := JWTClaims{
			user.ID,
			time.Now().Add(48 * time.Hour),
			jwt.StandardClaims{},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString([]byte(globalConfig.SigningSecret))
		if err != nil {
			log.Fatal(err)
		}
		result.Success = true
		result.Token = ss
	}
}

// List a users' tasks
func listTasks(w http.ResponseWriter, r *http.Request) {
	data := getTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Add a new task
func addTask(w http.ResponseWriter, r *http.Request) {
	type newTask struct {
		TaskName string `json:"taskName"`
	}
	var payload newTask
	json.NewDecoder(r.Body).Decode(&payload)
	w.Header().Set("Content-Type", "application/json")
	if err := insertTask(payload.TaskName); err != nil {
		json.NewEncoder(w).Encode(false)
		return
	}
	json.NewEncoder(w).Encode(true)
}

// Update completed column for a given task
func updateTaskStatus(w http.ResponseWriter, r *http.Request) {
	type taskID struct {
		TaskID    int64 `json:"taskID"`
		Completed bool  `json:"completed"`
	}
	var payload taskID
	json.NewDecoder(r.Body).Decode(&payload)
	w.Header().Set("Content-Type", "application/json")
	if err := updateTask(payload.TaskID, payload.Completed); err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(false)
		return
	}
	json.NewEncoder(w).Encode(true)
}

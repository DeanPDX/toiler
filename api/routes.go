package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type credentials struct {
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

	var creds credentials
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
		writeError(w, err.Error())
		return
	}
	insertUser(creds.Email, hash)
	user = getUserByEmail(creds.Email)
	// Generate our token
	signedString, err := generateToken(user.ID)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	result.Success = true
	result.Token = signedString
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := authResult{Success: false}
	defer json.NewEncoder(w).Encode(&result)

	var creds credentials
	json.NewDecoder(r.Body).Decode(&creds)
	user := getUserByEmail(creds.Email)
	if user.ID == 0 {
		return
	}
	if CheckPasswordHash(creds.Password, user.PasswordHash) {
		// Generate our token
		signedString, err := generateToken(user.ID)
		if err != nil {
			writeError(w, err.Error())
			return
		}
		result.Success = true
		result.Token = signedString
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	items := ""
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			items += fmt.Sprintf("<li>%v</li>", path)
			return nil
		})
	if err != nil {
		writeError(w, err.Error())
		return
	}

	healthCheckHTML := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	</head>
	<body>
	  <h1>App Health Check</h1>
	  <p>The app is working</p>
	  <p>Files:</p>
	  <ul>
	  %v
	  </ul>
	</body>
	</html>`, items)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(healthCheckHTML))
}

// List a users' tasks
func listTasks(w http.ResponseWriter, r *http.Request) {
	jwt, err := parseToken(r.Header.Get("X-Auth"))
	if err != nil {
		writeError(w, err.Error())
		return
	}
	data := getTasks(jwt.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Add a new task
func addTask(w http.ResponseWriter, r *http.Request) {
	jwt, err := parseToken(r.Header.Get("X-Auth"))
	if err != nil {
		writeError(w, err.Error())
		return
	}
	type newTask struct {
		TaskName string `json:"taskName"`
	}
	var payload newTask
	json.NewDecoder(r.Body).Decode(&payload)
	w.Header().Set("Content-Type", "application/json")

	if err := insertTask(payload.TaskName, jwt.UserID); err != nil {
		writeError(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(true)
}

// Update completed column for a given task
func updateTaskStatus(w http.ResponseWriter, r *http.Request) {
	jwt, err := parseToken(r.Header.Get("X-Auth"))
	if err != nil {
		writeError(w, err.Error())
		return
	}
	type taskID struct {
		TaskID    int64 `json:"taskID"`
		Completed bool  `json:"completed"`
	}
	var payload taskID
	json.NewDecoder(r.Body).Decode(&payload)
	w.Header().Set("Content-Type", "application/json")
	if err := updateTask(payload.TaskID, jwt.UserID, payload.Completed); err != nil {
		writeError(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(true)
}

func writeError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errorMessage))
}

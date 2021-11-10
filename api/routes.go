package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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

func healthCheck(w http.ResponseWriter, r *http.Request) {
	items := ""
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			items += fmt.Sprintf("<li>%v</li>", path)
			//fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
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

func healthCheckJSON(w http.ResponseWriter, r *http.Request) {
	items := make([]string, 0)
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			items = append(items, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// List a users' tasks
func listTasks(w http.ResponseWriter, r *http.Request) {
	// TODO: handle errors
	jwt, err := parseToken(r.Header.Get("X-Auth"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(jwt)
	data := getTasks(jwt.UserID)
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
	jwt, err := parseToken(r.Header.Get("X-Auth"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err := insertTask(payload.TaskName, jwt.UserID); err != nil {
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

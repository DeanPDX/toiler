package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Handle initialization
	setupConfig()
	initializeDB(globalConfig.DSN)
	initializeJWT(globalConfig.SigningSecret)
	defer closeDB()

	// Set up routes
	http.HandleFunc("/api/createAccount", createAccount)
	http.HandleFunc("/api/authenticate", authenticate)
	http.HandleFunc("/health", healthCheck)
	http.Handle("/api/tasks/list", mustBeAuthenticated(http.HandlerFunc(listTasks)))
	http.Handle("/api/tasks/excel", mustBeAuthenticated(http.HandlerFunc(listTasksExcel)))
	http.Handle("/api/tasks/add", mustBeAuthenticated(http.HandlerFunc(addTask)))
	http.Handle("/api/tasks/update", mustBeAuthenticated(http.HandlerFunc(updateTaskStatus)))
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Listen and serve content
	fmt.Println("\nListening on port", globalConfig.Port)
	log.Fatal(http.ListenAndServe(":"+globalConfig.Port, nil))
}

func mustBeAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth") == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

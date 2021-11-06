package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// AppConfig stores our app's configuration.
type AppConfig struct {
	Database      string `json:"database"`
	SigningSecret string `json:"signingSecret"`
	Port          string `json:"port"`
}

var globalConfig AppConfig

// Read config and set global config var
func readConfig() {
	fmt.Println("Reading config file")
	var config AppConfig
	// Open our config and read it.
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error trying to read config.json:", err)
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&config)
	require(config.Database, "Database")
	require(config.SigningSecret, "JWT signing secret")
	// Default port to 8090.
	if config.Port == "" {
		config.Port = ":8090"
	}
	globalConfig = config
}

// require is a shortcut for making sure our config contains a value and log.fatal-ing in the the event it doesn't.
func require(value, valueName string) {
	if value == "" {
		log.Fatalf("%v is required in config file", valueName)
	}
}

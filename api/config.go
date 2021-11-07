package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// AppConfig stores our app's configuration.
type AppConfig struct {
	DSN           string `json:"dsn"`
	SigningSecret string `json:"signingSecret"`
	Port          string `json:"port"`
}

var globalConfig AppConfig

// Set up config and set global config var
func setupConfig() {
	globalConfig = readConfigFile()
}

// readConfigFile attempts to read a config file. Probably for use during development. If no file is found, returns default values.
func readConfigFile() AppConfig {
	fmt.Println("Reading config file")
	var config AppConfig
	// Open our config and read it.
	f, err := os.Open("config.json")
	if err != nil {
		log.Println("conf.json is not present. Using environment variables.")
		return getEnvironmentVars()
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&config)
	// If we do have a config file, make sure it has required items.
	validateConfig(config, "conf.json")
	return config
}

func getEnvironmentVars() AppConfig {
	fmt.Println("Reading environment variables")
	// Create struct and populate it with environment variables.
	config := AppConfig{}
	config.DSN = os.Getenv("DSN")
	config.Port = os.Getenv("PORT")
	config.SigningSecret = os.Getenv("SIGNINGSECRET")
	fmt.Println("Config:", config)
	// Make sure our environment variables contain all required config settings.
	validateConfig(config, "environment variables")
	return config
}

// validateConfig will make sure all required fields are present in config, regardless of source. sourceName should be a user-readable string indicating the source of the config.
func validateConfig(config AppConfig, sourceName string) {
	require(config.DSN, "DSN", sourceName)
	require(config.SigningSecret, "JWT signing secret", sourceName)
	require(config.Port, "Port", sourceName)
}

// require is a shortcut for making sure our config contains a value and log.fatal-ing in the the event it doesn't.
func require(value, valueName, source string) {
	if value == "" {
		log.Fatalf("%v is required in %v", valueName, source)
	}
}

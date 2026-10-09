package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/clyvecute/configra/internal/config"
	"github.com/clyvecute/configra/internal/configs"
	"github.com/clyvecute/configra/internal/db"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	schemaPath := validateCmd.String("schema", "schema.json", "Path to the schema file")
	configPath := validateCmd.String("config", "config.json", "Path to the configuration file")

	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	pushFile := pushCmd.String("file", "config.json", "Config file to push")
	pushProject := pushCmd.String("project", "", "Project ID")
	pushKey := pushCmd.String("key", "feature_flags", "Config key name to push under")
	pushEnv := pushCmd.String("env", "1", "Environment ID")
	pushHost := pushCmd.String("host", "http://localhost:8080", "API Host URL")

	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchProject := fetchCmd.String("project", "", "Project ID")
	fetchKey := fetchCmd.String("key", "", "Config key to fetch")
	fetchEnv := fetchCmd.String("env", "1", "Environment ID")
	fetchHost := fetchCmd.String("host", "http://localhost:8080", "API Host URL")

	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
	rollbackProject := rollbackCmd.String("project", "", "Project ID")
	rollbackKey := rollbackCmd.String("key", "", "Config Key")
	rollbackVersion := rollbackCmd.String("version", "", "Target Version to restore")
	rollbackHost := rollbackCmd.String("host", "http://localhost:8080", "API Host URL")
	rollbackEnv := rollbackCmd.String("env", "1", "Environment ID")

	switch os.Args[1] {
	case "validate":
		validateCmd.Parse(os.Args[2:])
		runValidate(*schemaPath, *configPath)
	case "push":
		pushCmd.Parse(os.Args[2:])
		runPush(*pushFile, *pushProject, *pushKey, *pushEnv, *pushHost)
	case "fetch":
		fetchCmd.Parse(os.Args[2:])
		runFetch(*fetchProject, *fetchKey, *fetchEnv, *fetchHost)
	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
		runRollback(*rollbackProject, *rollbackKey, *rollbackVersion, *rollbackEnv, *rollbackHost)
	case "migrate":
		// Ensure we load config to get DB creds
		runMigrate()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Configra CLI")
	fmt.Println("Usage:")
	fmt.Println("  validate -schema <path> -config <path>                          Validate a config against a schema locally")
	fmt.Println("  push     -file <path> -project <id> -key <key> -env <env_id>   Push a config to the server")
	fmt.Println("  fetch    -project <id> -key <key> -env <env_id>                Fetch active config from server")
	fmt.Println("  rollback -project <id> -key <key> -version <n>                 Roll back to a specific version")
	fmt.Println("  migrate                                                          Run database migrations")
}

func runMigrate() {
	fmt.Println("Running migrations...")
	cfg := config.Load()
	database, err := db.Connect(cfg.DB)
	if err != nil {
		fmt.Printf("Failed to connect to DB: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	cwd, _ := os.Getwd()
	// Assumption: running from project root or having migrations folder relative
	migrationsDir := filepath.Join(cwd, "internal", "db", "migrations")

	if err := db.Migrate(database, migrationsDir); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Migrations completed successfully.")
}

func runValidate(schemaFile, configFile string) {
	fmt.Printf("Validating %s against %s...\n", configFile, schemaFile)

	// Read Schema
	sBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	var schema configs.Schema
	if err := json.Unmarshal(sBytes, &schema); err != nil {
		fmt.Printf("Error parsing schema JSON: %v\n", err)
		os.Exit(1)
	}

	// Read Config
	cBytes, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(cBytes, &config); err != nil {
		fmt.Printf("Error parsing config JSON: %v\n", err)
		os.Exit(1)
	}

	// Validate
	if err := configs.Validate(schema, config); err != nil {
		fmt.Printf("\u274C Validation FAILED: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\u2705 Configuration is VALID.")
}

func runPush(configFile, projectID, cfgKey, envIDStr, host string) {
	schemaFile := "schema.json"

	cBytes, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	sBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Printf("Error reading schema: %v\n", err)
		os.Exit(1)
	}

	var configMap map[string]interface{}
	var schemaMap map[string]interface{}
	json.Unmarshal(cBytes, &configMap)
	json.Unmarshal(sBytes, &schemaMap)

	// Validate locally first
	var schemaStruct configs.Schema
	json.Unmarshal(sBytes, &schemaStruct)
	if err := configs.Validate(schemaStruct, configMap); err != nil {
		fmt.Printf("Validation failed locally: %v\n", err)
		os.Exit(1)
	}

	pID, err := strconv.Atoi(projectID)
	if err != nil || pID <= 0 {
		fmt.Println("A valid positive -project ID is required")
		os.Exit(1)
	}
	eID, err := strconv.Atoi(envIDStr)
	if err != nil || eID <= 0 {
		fmt.Println("A valid positive -env ID is required")
		os.Exit(1)
	}
	if cfgKey == "" {
		fmt.Println("-key is required")
		os.Exit(1)
	}
	apiKey := os.Getenv("CONFIGRA_API_KEY")
	if apiKey == "" {
		fmt.Println("CONFIGRA_API_KEY environment variable is not set")
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"project_id": pID,
		"env_id":     eID,
		"key":        cfgKey,
		"data":       configMap,
		"schema":     schemaMap,
	}

	body, _ := json.Marshal(payload)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/configs", host), bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to build request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", apiKey)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		fmt.Printf("Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned error: %s\n", resp.Status)
		return
	}

	fmt.Println("Successfully pushed config to server!")
}

func runFetch(projectID, key, envIDStr, host string) {
	pID, err := strconv.Atoi(projectID)
	if err != nil || pID <= 0 {
		fmt.Println("A valid positive -project ID is required")
		os.Exit(1)
	}
	eID, err := strconv.Atoi(envIDStr)
	if err != nil || eID <= 0 {
		fmt.Println("A valid positive -env ID is required")
		os.Exit(1)
	}
	if key == "" {
		fmt.Println("-key is required")
		os.Exit(1)
	}

	apiKey := os.Getenv("CONFIGRA_API_KEY")
	if apiKey == "" {
		fmt.Println("CONFIGRA_API_KEY environment variable is not set")
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/v1/configs/%s?env_id=%d", host, key, eID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Failed to build request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Fetch failed: %s\n", resp.Status)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		os.Exit(1)
	}

	formatted, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(formatted))
}

func runRollback(projectID, key, version, env, host string) {
	// Simple conversions
	pID, err := strconv.Atoi(projectID)
	if err != nil || pID <= 0 {
		fmt.Println("A valid positive -project ID is required")
		os.Exit(1)
	}
	vID := 0
	if n, e := strconv.Atoi(version); e == nil {
		vID = n
	}
	eID, err := strconv.Atoi(env)
	if err != nil || eID <= 0 {
		fmt.Println("A valid positive -env ID is required")
		os.Exit(1)
	}
	if key == "" || vID <= 0 {
		fmt.Println("-key and a positive -version are required")
		os.Exit(1)
	}
	apiKey := os.Getenv("CONFIGRA_API_KEY")
	if apiKey == "" {
		fmt.Println("CONFIGRA_API_KEY environment variable is not set")
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"project_id":     pID,
		"env_id":         eID,
		"key":            key,
		"target_version": vID,
	}

	body, _ := json.Marshal(payload)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/configs/%s/rollback", host, key), bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to build request: %v\n", err)
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", apiKey)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		fmt.Printf("Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Rollback failed: %s\n", resp.Status)
		return
	}

	fmt.Printf("\u2705 Successfully rolled back '%s' to version %s!\n", key, version)
}

// Add these imports at the top if missing: bytes, net/http

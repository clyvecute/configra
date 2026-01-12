package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

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
	_ = pushCmd.String("file", "config.json", "Config file to push")
	_ = pushCmd.String("project", "", "Project ID")

	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	_ = fetchCmd.String("project", "", "Project ID")
	_ = fetchCmd.String("env", "prod", "Environment")

	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
	_ = rollbackCmd.String("project", "", "Project ID")
	_ = rollbackCmd.String("key", "", "Config Key")
	_ = rollbackCmd.String("version", "", "Target Version to restore")

	switch os.Args[1] {
	case "validate":
		validateCmd.Parse(os.Args[2:])
		runValidate(*schemaPath, *configPath)
	case "push":
		pushCmd.Parse(os.Args[2:])
		file := "config.json"
		if f := pushCmd.Lookup("file"); f != nil {
			file = f.Value.String()
		}
		proj := ""
		if p := pushCmd.Lookup("project"); p != nil {
			proj = p.Value.String()
		}
		runPush(file, proj)
	case "fetch":
		fetchCmd.Parse(os.Args[2:])
		fmt.Println("Fetch logic to be implemented. Would fetch from API.")
	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
		p := ""; if f := rollbackCmd.Lookup("project"); f != nil { p = f.Value.String() }
		k := ""; if f := rollbackCmd.Lookup("key"); f != nil { k = f.Value.String() }
		v := ""; if f := rollbackCmd.Lookup("version"); f != nil { v = f.Value.String() }
		runRollback(p, k, v)
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
	fmt.Println("  validate -schema <path> -config <path>   Validate a config against a schema locally")
	fmt.Println("  push     -file <path> -project <id>      Push a config to the server")
	fmt.Println("  fetch    -project <id> -env <name>       Fetch active config from server")
	fmt.Println("  migrate                                  Run database migrations")
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

func runPush(configFile, projectID string) {
	// 1. Read the config file and assumed schema file (for now co-located or we should bundle them)
	// For this demo, let's assume schema.json is in the same dir
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
	
	// 2. Validate locally first
	var schemaStruct configs.Schema
	json.Unmarshal(sBytes, &schemaStruct)
	if err := configs.Validate(schemaStruct, configMap); err != nil {
		fmt.Printf("Validation failed locally: %v\n", err)
		os.Exit(1)
	}

	// 3. Send to API
	// Construct payload matching CreateRequest in handler
	// We need to parse projectID to int, let's default to 1 if empty for push demo
	pID := 1 // Default
	eID := 1 // Default env
	
	payload := map[string]interface{}{
		"project_id": pID,
		"env_id":     eID,
		"key":        "feature_flags", // Default key for now
		"data":       configMap,
		"schema":     schemaMap,
	}
	
	body, _ := json.Marshal(payload)
	
	resp, err := http.Post("http://localhost:8080/v1/configs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned error: %s\n", resp.Status)
		// Read body for details
		return
	}

	fmt.Println("Successfully pushed config to server!")
}

func runRollback(projectID, key, version string) {
	// Simple conversions
	pID := 1 // default
	vID := 0
	fmt.Sscanf(version, "%d", &vID)

	payload := map[string]interface{}{
		"project_id":     pID,
		"env_id":         1, // default
		"key":            key,
		"target_version": vID,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8080/v1/rollback", "application/json", bytes.NewBuffer(body))
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


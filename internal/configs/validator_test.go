package configs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	// Define a schema
	schemaJSON := `
	{
		"version": 1,
		"rules": {
			"feature_enabled": { "type": "bool", "required": true },
			"max_users": { "type": "int", "min": 1, "max": 100 },
			"region": { "type": "enum", "allowed": ["us-east", "eu-west"] }
		}
	}`

	var schema Schema
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			name: "Valid Config",
			config: `{
				"feature_enabled": true,
				"max_users": 50,
				"region": "us-east"
			}`,
			wantErr: false,
		},
		{
			name: "Missing Required Field",
			config: `{
				"max_users": 50
			}`,
			wantErr: true,
		},
		{
			name: "Invalid Type",
			config: `{
				"feature_enabled": "yes", 
				"max_users": 50
			}`,
			wantErr: true,
		},
		{
			name: "Constraint Violation (Max)",
			config: `{
				"feature_enabled": true,
				"max_users": 150
			}`,
			wantErr: true,
		},
		{
			name: "Enum Violation",
			config: `{
				"feature_enabled": true,
				"region": "ap-south"
			}`,
			wantErr: true,
		},
		{
			name: "Unknown Field",
			config: `{
				"feature_enabled": true,
				"random_field": 123
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(tt.config), &config); err != nil {
				t.Fatalf("failed to parse config json: %v", err)
			}

			err := Validate(schema, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDefaultsAndBoundaries(t *testing.T) {
	min, max := float64(1), float64(10)
	schema := Schema{Rules: map[string]FieldRule{
		"count": {Type: TypeInt, Required: true, Min: &min, Max: &max},
		"label": {Type: TypeString, Default: "default"},
	}}
	config := map[string]interface{}{"count": float64(10)}
	if err := Validate(schema, config); err != nil {
		t.Fatalf("valid boundary config: %v", err)
	}
	if config["label"] != "default" {
		t.Fatalf("default not applied: %#v", config)
	}

	for _, value := range []interface{}{float64(0), float64(11), float64(1.5)} {
		t.Run(fmt.Sprintf("invalid_%v", value), func(t *testing.T) {
			got := map[string]interface{}{"count": value}
			if err := Validate(schema, got); err == nil {
				t.Fatalf("Validate(%v) unexpectedly succeeded", value)
			}
		})
	}
}

func TestValidateRejectsUnknownAndInvalidEnum(t *testing.T) {
	schema := Schema{Rules: map[string]FieldRule{"mode": {Type: TypeEnum, Allowed: []interface{}{"safe", "fast"}}}}
	for name, config := range map[string]map[string]interface{}{
		"unknown": {"mode": "safe", "extra": true},
		"enum":    {"mode": "unsafe"},
	} {
		t.Run(name, func(t *testing.T) {
			if err := Validate(schema, config); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

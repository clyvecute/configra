package configs

import (
	"encoding/json"
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

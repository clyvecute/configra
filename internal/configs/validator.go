package configs

import (
	"fmt"
	"reflect"
	"strings"
)

// DataType defines the supported types for configuration values.
type DataType string

const (
	TypeString  DataType = "string"
	TypeInt     DataType = "int"
	TypeFloat   DataType = "float"
	TypeBool    DataType = "bool"
	TypeEnum    DataType = "enum"
	TypeJSON    DataType = "json" // For complex nested objects
)

// FieldRule defines the validation logic for a single configuration key.
type FieldRule struct {
	Type        DataType      `json:"type"`
	Required    bool          `json:"required"`
	Description string        `json:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Min         *float64      `json:"min,omitempty"`  // For int/float
	Max         *float64      `json:"max,omitempty"`  // For int/float
	Allowed     []interface{} `json:"allowed,omitempty"` // For enum
}

// Schema defines the contract that a configuration must adhere to.
type Schema struct {
	Version int                  `json:"version"`
	Rules   map[string]FieldRule `json:"rules"`
}

// ValidationError represents a collection of validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Errors, "; "))
}

// Validate checks a raw configuration map against the provided Schema.
func Validate(schema Schema, config map[string]interface{}) error {
	var errs []string

	for key, rule := range schema.Rules {
		// 1. Check for missing required fields
		val, exists := config[key]
		if !exists {
			if rule.Required {
				errs = append(errs, fmt.Sprintf("field '%s' is required", key))
			} else if rule.Default != nil {
				// Apply default if missing (handled by caller usually, but good to know)
				config[key] = rule.Default
			}
			continue
		}

		// 2. Type validation
		if !isValidType(val, rule.Type) {
			errs = append(errs, fmt.Sprintf("field '%s' expected type %s, got %T", key, rule.Type, val))
			continue
		}

		// 3. Constraint validation
		if err := validateConstraints(key, val, rule); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// 4. Check for unknown fields (Strict mode)
	for key := range config {
		if _, known := schema.Rules[key]; !known {
			errs = append(errs, fmt.Sprintf("unknown field '%s' is not allowed by schema", key))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}

	return nil
}

// isValidType checks if the value matches the expected DataType using reflection.
// Note: JSON decoding often treats numbers as float64.
func isValidType(val interface{}, expected DataType) bool {
	switch expected {
	case TypeString:
		_, ok := val.(string)
		return ok
	case TypeBool:
		_, ok := val.(bool)
		return ok
	case TypeInt:
		// JSON unmarshal makes numbers float64, so we check if it's convertible to int perfectly
		f, ok := val.(float64)
		if ok {
			return f == float64(int(f))
		}
		// Also handle explicit int if passed from code
		_, ok = val.(int)
		return ok
	case TypeFloat:
		_, ok := val.(float64)
		return ok
	case TypeEnum:
		// Enum base type must be string or int generally, but we'll check against allowed list later.
		// For basic type check, we treat enums as comparable values (string/number).
		switch val.(type) {
		case string, float64, int, bool:
			return true
		}
		return false
	case TypeJSON:
		// Map or slice
		kind := reflect.TypeOf(val).Kind()
		return kind == reflect.Map || kind == reflect.Slice
	default:
		return false
	}
}

func validateConstraints(key string, val interface{}, rule FieldRule) error {
	// Min/Max for numbers
	if rule.Min != nil || rule.Max != nil {
		numVal, isNum := toFloat(val)
		if isNum {
			if rule.Min != nil && numVal < *rule.Min {
				return fmt.Errorf("field '%s' must be >= %v", key, *rule.Min)
			}
			if rule.Max != nil && numVal > *rule.Max {
				return fmt.Errorf("field '%s' must be <= %v", key, *rule.Max)
			}
		}
	}

	// Allowed values (Enum)
	if len(rule.Allowed) > 0 {
		found := false
		for _, allowed := range rule.Allowed {
			if reflect.DeepEqual(val, allowed) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("field '%s' has invalid value '%v'; allowed: %v", key, val, rule.Allowed)
		}
	}

	return nil
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

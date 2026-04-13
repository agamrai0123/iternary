package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator validates data against schemas
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateField validates a single field value against field schema
func (v *Validator) ValidateField(fieldName string, value interface{}, schema *FieldSchema) *Error {
	// Check if required
	if schema.Required && isEmpty(value) {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("%s is required", fieldName),
		}
	}

	// If not required and empty, validation passes
	if !schema.Required && isEmpty(value) {
		return nil
	}

	// Check type and constraints
	switch schema.Type {
	case "email":
		if err := validateEmail(value); err != nil {
			return &Error{
				Field:   fieldName,
				Message: fmt.Sprintf("Invalid email format: %v", err),
			}
		}
	case "password":
		if err := validatePassword(value, schema); err != nil {
			return &Error{
				Field:   fieldName,
				Message: err.Error(),
			}
		}
	case "uuid":
		if err := validateUUID(value); err != nil {
			return &Error{
				Field:   fieldName,
				Message: fmt.Sprintf("Invalid UUID format: %v", err),
			}
		}
	case "string":
		if err := validateString(value, schema); err != nil {
			return &Error{
				Field:   fieldName,
				Message: err.Error(),
			}
		}
	case "number":
		if err := validateNumber(value); err != nil {
			return &Error{
				Field:   fieldName,
				Message: fmt.Sprintf("Invalid number: %v", err),
			}
		}
	}

	// Check enum if specified
	if len(schema.Enum) > 0 {
		if err := validateEnum(value, schema.Enum); err != nil {
			return &Error{
				Field:   fieldName,
				Message: fmt.Sprintf("Value must be one of: %s", strings.Join(schema.Enum, ", ")),
			}
		}
	}

	// Check pattern if specified
	if schema.Pattern != "" {
		if err := validatePattern(value, schema.Pattern); err != nil {
			return &Error{
				Field:   fieldName,
				Message: fmt.Sprintf("Does not match pattern: %s", schema.Pattern),
			}
		}
	}

	return nil
}

// ValidateObject validates entire object against schema
func (v *Validator) ValidateObject(data map[string]interface{}, schema *Schema) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []Error{},
	}

	// Check all fields in schema
	for fieldName, fieldSchema := range schema.Fields {
		value, exists := data[fieldName]
		if !exists && !fieldSchema.Required {
			continue
		}

		if !exists && fieldSchema.Required {
			result.Valid = false
			result.Errors = append(result.Errors, Error{
				Field:   fieldName,
				Message: fmt.Sprintf("%s is required", fieldName),
			})
			continue
		}

		// Validate field
		if err := v.ValidateField(fieldName, value, fieldSchema); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, *err)
		}
	}

	return result
}

// Helper validation functions

func isEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case nil:
		return true
	default:
		return false
	}
}

func validateEmail(value interface{}) error {
	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	// Simple email regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil || !matched {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func validatePassword(value interface{}, schema *FieldSchema) error {
	password, ok := value.(string)
	if !ok {
		return fmt.Errorf("password must be a string")
	}

	if len(password) < schema.MinLen {
		return fmt.Errorf("password must be at least %d characters", schema.MinLen)
	}

	if len(password) > schema.MaxLen {
		return fmt.Errorf("password must be no more than %d characters", schema.MaxLen)
	}

	return nil
}

func validateString(value interface{}, schema *FieldSchema) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if schema.MinLen > 0 && len(str) < schema.MinLen {
		return fmt.Errorf("must be at least %d characters", schema.MinLen)
	}

	if schema.MaxLen > 0 && len(str) > schema.MaxLen {
		return fmt.Errorf("must be no more than %d characters", schema.MaxLen)
	}

	return nil
}

func validateNumber(value interface{}) error {
	// Accept int, int64, float64
	switch value.(type) {
	case int, int64, float64:
		return nil
	default:
		return fmt.Errorf("must be a number")
	}
}

func validateUUID(value interface{}) error {
	uuid, ok := value.(string)
	if !ok {
		return fmt.Errorf("uuid must be a string")
	}

	// Simple UUID validation (v4 format)
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(uuidRegex, uuid)
	if err != nil || !matched {
		return fmt.Errorf("invalid UUID v4 format")
	}

	return nil
}

func validateEnum(value interface{}, validValues []string) error {
	strValue := fmt.Sprintf("%v", value)

	for _, valid := range validValues {
		if strValue == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid enum value")
}

func validatePattern(value interface{}, pattern string) error {
	strValue := fmt.Sprintf("%v", value)

	matched, err := regexp.MatchString(pattern, strValue)
	if err != nil || !matched {
		return fmt.Errorf("pattern mismatch")
	}

	return nil
}

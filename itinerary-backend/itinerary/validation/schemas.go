package validation

// FieldSchema defines validation rules for a single field
type FieldSchema struct {
	Required bool     `json:"required"`
	Type     string   `json:"type"` // "string", "email", "number", "uuid", "password"
	MinLen   int      `json:"min_len,omitempty"`
	MaxLen   int      `json:"max_len,omitempty"`
	Pattern  string   `json:"pattern,omitempty"` // regex
	Enum     []string `json:"enum,omitempty"`
}

// Schema defines validation rules for an entire request
type Schema struct {
	Fields map[string]*FieldSchema
}

// Predefined validation schemas

var (
	// User registration schema
	UserRegistrationSchema = &Schema{
		Fields: map[string]*FieldSchema{
			"username": {
				Required: true,
				Type:     "string",
				MinLen:   3,
				MaxLen:   50,
			},
			"email": {
				Required: true,
				Type:     "email",
			},
			"password": {
				Required: true,
				Type:     "password",
				MinLen:   8,
				MaxLen:   72,
			},
			"full_name": {
				Required: false,
				Type:     "string",
				MaxLen:   100,
			},
		},
	}

	// Login schema
	LoginSchema = &Schema{
		Fields: map[string]*FieldSchema{
			"email": {
				Required: true,
				Type:     "email",
			},
			"password": {
				Required: true,
				Type:     "password",
			},
		},
	}

	// MFA verify schema
	MFAVerifySchema = &Schema{
		Fields: map[string]*FieldSchema{
			"code": {
				Required: true,
				Type:     "string",
				MinLen:   6,
				MaxLen:   6,
			},
		},
	}

	// Link account schema
	LinkAccountSchema = &Schema{
		Fields: map[string]*FieldSchema{
			"code": {
				Required: true,
				Type:     "string",
			},
			"provider": {
				Required: true,
				Type:     "string",
				Enum:     []string{"github", "google", "microsoft"},
			},
		},
	}

	// Disable MFA schema
	DisableMFASchema = &Schema{
		Fields: map[string]*FieldSchema{
			"password": {
				Required: true,
				Type:     "password",
			},
		},
	}
)

// Error represents a validation error
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid  bool    `json:"valid"`
	Errors []Error `json:"errors,omitempty"`
}

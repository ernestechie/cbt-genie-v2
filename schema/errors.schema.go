package schema

// Custom error messages map
var CustomErrorMessages = map[string]string{
	"FirstName.required":  "First name is required",
	"FirstName.min":       "First name must be at least 3 characters long",
	"FirstName.max":       "First name cannot exceed 32 characters",
	"LastName.required":  "Last name is required",
	"LastName.min":       "Last name must be at least 3 characters long",
	"LastName.max":       "Last name cannot exceed 32 characters",
	"Email.required": "Email is required",
	"Email.email":    "Email must be a valid email address",
	"Age.gte":        "Age must be 12 or greater",
	"Age.lte":        "Age must be 100 or less",
}


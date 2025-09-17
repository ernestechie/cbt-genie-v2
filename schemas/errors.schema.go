package schema

// Custom error messages map
var CustomErrorMessages = map[string]string{
	"first_name.required":  "First name is required",
	"first_name.min":       "First name must be at least 3 characters long",
	"first_name.max":       "First name cannot exceed 32 characters",
	"last_name.required":  "Last name is required",
	"last_name.min":       "Last name must be at least 3 characters long",
	"last_name.max":       "Last name cannot exceed 32 characters",
	"email.required": "Email is required",
	"email.email":    "Email must be a valid email address",
	"age.gte":        "Age must be 12 or greater",
	"age.lte":        "Age must be 100 or less",
	"otp_code.required": "Invalid OTP",
	"otp_code.min": "Invalid OTP",
	"otp_code.max": "Invalid OTP",
	"otp_token.required": "Invalid OTP token",	
	"otp_token.min": "Invalid OTP token",	
}

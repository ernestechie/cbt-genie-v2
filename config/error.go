package config

// ErrorResponse for custom error formatting
type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

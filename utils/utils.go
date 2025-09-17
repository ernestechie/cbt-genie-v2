package utils

import (
	"crypto/rand"
	"log"
	"os"
	"strings"

	"github.com/ernestechie/cbt-genie-v2/config"
	schemas "github.com/ernestechie/cbt-genie-v2/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

// Initialize validator
var Validate = validator.New()

// Hash a plain text string and return the hashed value

const otpChars = "1234567890"

func GenerateSecureOtp(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for index := range length {
		buffer[index] = otpChars[int(buffer[index])%otpCharsLength]
	}

	return string(buffer), nil
}

// Hash a plain text string and return the hashed value
func GetHashedValue(inputVal string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(inputVal), 12)
	return string(bytes), err
}

// Validate plain text string with corresponding hashend value, return "true" or "false"
func ValidateHashedValue(inputVal, hashedValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(inputVal))
	return err == nil
}

// Parses the request body into the provided struct and validates it
func ParseAndValidate(c *fiber.Ctx, data any) []config.ErrorResponse {
	// Parse JSON body into struct
	if err := c.BodyParser(data); err != nil {
		log.Printf("Error_parsing_body: %+v", err)
		return []config.ErrorResponse{
			{
				Field: "",
				Error: "Invalid request body format",
			},
		}
	}

	// Initialize validator
	validate := validator.New()

	// Validate struct
	if err := validate.Struct(data); err != nil {
		log.Println("Validation_error \n", err)

		var errors []config.ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			// Create key for custom error message (e.g., "Name.required")
			key := strings.ToLower(err.Field() + "." + err.Tag())
			// Get custom message or fallback to default
			message, exists := schemas.CustomErrorMessages[key]
			if !exists {
				message = err.Error()
			}
			errors = append(errors, config.ErrorResponse{
				Field: err.Field(),
				Error: message,
			})
		}
		return errors
	}

	return nil
}

func VerifyJwt(tokenString string) (jwt.MapClaims, error) {
	var jwtSecret string
	if jwtSecret = os.Getenv("JWT_SECRET"); jwtSecret == "" {
		log.Fatal("You must set your 'JWT_SECRET' environment variable.")
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		hmacSampleSecret := []byte(jwtSecret)
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

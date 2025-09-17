package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ernestechie/cbt-genie-v2/config"
	"github.com/ernestechie/cbt-genie-v2/models"
	schemas "github.com/ernestechie/cbt-genie-v2/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

// Initialize validator
var Validate = validator.New()
// Set to 5 minutes

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

type JwtUser struct {
	TokenExpiry 	time.Time
	UserId				string
}

func VerifyJwt(tokenString string) (*JwtUser, error) {
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
	
	fmt.Println("VerifyJwt (tokenString):", tokenString)
	if err != nil {
		return &JwtUser{}, err
	}

	claims := token.Claims.(jwt.MapClaims)
	sub :=  claims["sub"]
	userId := fmt.Sprintf("%v", sub)
	if userId == "" {
		return &JwtUser{}, errors.New("invalid token. cannot find user")
	} 

	var expTime time.Time
	exp :=	claims["exp"]
	// Using type assertion
	if expFloat, ok := exp.(float64); ok {
		expTime = time.Unix(int64(expFloat), 0)
	}

	return &JwtUser{
		UserId: userId,
		TokenExpiry: expTime,
	}, nil
}

type NewJwtUserReturn struct {
	RefreshToken 	string
	AccessToken		string
}

func SignJwt(user *models.User) (*NewJwtUserReturn, error) {
	var jwtSecret string
	if jwtSecret = os.Getenv("JWT_SECRET"); jwtSecret == "" {
		log.Fatal("You must set your 'JWT_SECRET' environment variable.")
	}

	// Create a new token object, specifying signing method and the claims you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 2).Unix(), // 2 minutes
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	accessTokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return &NewJwtUserReturn{}, err
	}

	// TODO: Move refresh token to another function, and use a different signing protocol, more secured.
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return &NewJwtUserReturn{}, err
	}

	return &NewJwtUserReturn{
		RefreshToken: refreshTokenString,
		AccessToken: accessTokenString,
	}, nil
}

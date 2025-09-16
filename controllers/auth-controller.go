package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ernestechie/cbt-genie-v2/models"
	"github.com/ernestechie/cbt-genie-v2/utils"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GET STARTED
type getStartedRequestBody struct {
	Email		string 		`json:"email" validate:"required,email"`
}

// GET STARTED - /auth/get-started 
func GetStarted (c *fiber.Ctx) error {
	var reqBody getStartedRequestBody
	var foundUser models.User
	var responseMessage string
	// var nextStep string

	if errs := utils.ParseAndValidate(c, &reqBody); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": errs,
			"message": "Error processing request",
		})
	}

	// Check if user exist with that email.
	filter := bson.D{{Key: "email", Value: strings.ToLower(reqBody.Email)}}

	// Generate a secure OTP
	otpCode, err  := utils.GenerateSecureOtp(6)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": err,
			"message": "Error! Cannot authenticate user at the moment.",
		})
	}

	log.Println("OTP_CODE ", otpCode)

	hashedOtpToken, err := utils.GetHashedValue(otpCode)
		if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": err,
			"message": "Error processing request",
		})
	}
	
	// OTP Expiry time is 5 minutes
	otpExpiry := time.Now().Local().Add(time.Minute * 5)

	// Check if user exists
	userExistsErr := usersColl.FindOne(context.TODO(), filter).Decode(&foundUser)
		if userExistsErr != nil {
			// TODO: FIND A BETTER WAY TO DO THIS
			if userExistsErr.Error() == "mongo: no documents in result" {
				// nextStep = "Onboarding"
				responseMessage = "Account created succesfully! An OTP has been sent to your email."

				log.Println("No existing. Creating new user...")

				// No user Found, build user object.
				foundUser.Email = strings.ToLower(reqBody.Email)
				foundUser.OtpToken = hashedOtpToken
				foundUser.OtpExpiry = otpExpiry

				foundUser.CreatedAt = time.Now()
				foundUser.UpdatedAt = time.Now()

				// Create new user if the user does not exist.
				res, err := usersColl.InsertOne(context.TODO(), foundUser)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"success": false,
						"errors": err.Error(),
						"message": "Error occured while creating user account",
					})
				}

				userId := res.InsertedID.(bson.ObjectID)
				foundUser.ID = userId
				} 
			} else {
				// nextStep = "Login"
				responseMessage = "An OTP has been sent to your email."

				update := bson.D{
					{Key: "$set", Value: bson.D{
							{Key: "otp_token", Value: hashedOtpToken},
							{Key: "otp_expiry", Value: otpExpiry},
					}},
				}

			opts := options.UpdateOne().SetUpsert(true)
			_, err := usersColl.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"errors": err.Error(),
					"message": "Error authenticating user",
				})
			}
		}

	// TODO: Send secure otp to user email.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": fiber.Map{
				"email": reqBody.Email,
				"id": foundUser.ID,
			},
		},
		"message": responseMessage,
	})
}

// VERIFY OTP
type verifyOtpRequestBody struct {
	Email			string 		`json:"email" validate:"required,email"`
	OtpCode		string 		`json:"otp_code" validate:"required,min=6,max=6"`
}

// GET STARTED - /auth/verify-otp
func VerifyOtp (c *fiber.Ctx) error {
	var reqBody verifyOtpRequestBody
	var foundUser models.User

	// Parse and validate request body using utils
	if errs := utils.ParseAndValidate(c, &reqBody); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": errs,
			"message": "Error processing request",
		})
	}

	// TODO: Check if user exists.
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: strings.ToLower(reqBody.Email)}}}},
			}},
	}

	// Check if user exists
	userExistsErr := usersColl.FindOne(context.TODO(), filter).Decode(&foundUser)
		if userExistsErr != nil { 
			// TODO: FIND A BETTER WAY TO DO THIS
			if userExistsErr.Error() == "mongo: no documents in result" {
				// nextStep = "Onboarding"
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"message": "User not found",
				})
			}
		} 
	
	otpIsValid := utils.ValidateHashedValue(reqBody.OtpCode, foundUser.OtpToken)
	otpExpired := time.Now().Local().After(foundUser.OtpExpiry)

	fmt.Println("Otp_is_valid: ", otpIsValid)
	fmt.Println("Otp_expired: ", time.Now().After(foundUser.OtpExpiry))
	// If there is an otp token in our user, or the otp token from request body is invalid, OR the token expiry is greater than now
	if len(foundUser.OtpToken) == 0 || (len(foundUser.OtpToken) > 1 && !otpIsValid  || otpExpired) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired OTP",
		})
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "otp_token", Value: ""}}}}
	opts := options.UpdateOne().SetUpsert(true)
	_, err := usersColl.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": err.Error(),
			"message": "Error authenticating user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": fiber.Map{
				"email": reqBody.Email,
				"id": foundUser.ID,
			},
		},
		"message": "Email verified successfully",
	})
}

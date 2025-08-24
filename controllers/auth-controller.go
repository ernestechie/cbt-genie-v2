package controllers

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ernestechie/cbt-genie-v2/models"
	"github.com/ernestechie/cbt-genie-v2/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func GetUserUpdates (user models.User, filter bson.D) (bson.D, error) {
	updates := bson.D{}

	// get the type of struct == User
	typeData := reflect.TypeOf(user)

	// get the values from the provided object: name -> John Smith
	values := reflect.ValueOf(user)

	// starting from index 1 to exclude the ID field
	for i := range typeData.NumField() {
		field := typeData.Field(i)   // get the field from the struct definition
		val := values.Field(i)       // get the value from the specified field position
		tag := field.Tag.Get("json") // from the field, get the json struct tag

		// we want to avoid zero values, as the omitted fields from newBook
		// corresponds to their zero values, and we only want provided fields
		if !val.IsZero() && tag != "id" {
			update := bson.E{Key: tag, Value: val.Interface()}
			updates = append(updates, update)
		}
	}

	updateFilter := bson.D{{Key: "$set", Value: updates}}
	// _, err := usersColl.UpdateOne(context.TODO(), filter, updateFilter)
	// if err != nil {
	// 	log.Fatalf("error updating user: %v", err)
	// 	return nil, err
	// }

	return updateFilter, nil
}


//*
// GET STARTED
// */

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
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: strings.ToLower(reqBody.Email)}}}},
			}},
	}

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

	foundUser.OtpToken = hashedOtpToken
	foundUser.OtpExpiry = otpExpiry

	// Check if user exists
	userExistsErr := usersColl.FindOne(context.TODO(), filter).Decode(&foundUser)
		if userExistsErr != nil {
			log.Print(userExistsErr)

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
					log.Println(err)
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
				// create updates variable to hold all the update fields
				// updates := bson.D{}
				updates, _ := GetUserUpdates(foundUser, filter)

				// nextStep = "Login"
				responseMessage = "An OTP has been sent to your email."

				// User Exists
				foundUser.OtpToken = hashedOtpToken
				foundUser.OtpExpiry = otpExpiry

				// get the type of struct == User
				// typeData := reflect.TypeOf(foundUser)

				// // get the values from the provided object: name -> John Smith
				// values := reflect.ValueOf(foundUser)

				// // starting from index 1 to exclude the ID field
				// for i := range typeData.NumField() {
				// 	field := typeData.Field(i)   // get the field from the struct definition
				// 	val := values.Field(i)       // get the value from the specified field position
				// 	tag := field.Tag.Get("json") // from the field, get the json struct tag

				// 	// we want to avoid zero values, as the omitted fields from newBook
				// 	// corresponds to their zero values, and we only want provided fields
				// 	if !val.IsZero() && tag != "id" {
				// 		update := bson.E{Key: tag, Value: val.Interface()}
				// 		updates = append(updates, update)
				// 	}
				// }

				updateFilter := bson.D{{Key: "$set", Value: updates}}
				_, err := usersColl.UpdateOne(context.TODO(), filter, updateFilter)
				if err != nil {
					log.Fatalf("error updating user: %v", err)
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
			"email": foundUser.Email,
			"id": foundUser.ID,
			// "metadata": fiber.Map{
			// 	"next_step":	nextStep,
			// 	"time_stamp": time.Now().Local(),
			// },
		},
		"message": responseMessage,
	})
}


//*
// VERIFY OTP
// */


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
	// If there is an otp token in our user, the otp token from request body is invalid, OR the token expiry is greater than now
	if len(foundUser.OtpToken) == 0 || (len(foundUser.OtpToken) > 1 && !otpIsValid  || time.Now().After(foundUser.OtpExpiry)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired OTP",
		})
	}

	foundUser.OtpToken = ""

	updates, _ := GetUserUpdates(foundUser, filter)

	updateFilter := bson.D{{Key: "$set", Value: updates}}
	_, err := usersColl.UpdateOne(context.TODO(), filter, updateFilter)
	if err != nil {
		log.Fatalf("error updating user: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error authenticating user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"email": reqBody.Email,
			"id": foundUser.ID,
		},
		"message": "Successful",
	})
}

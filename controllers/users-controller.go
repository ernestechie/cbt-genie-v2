package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ernestechie/cbt-genie-v2/db"
	"github.com/ernestechie/cbt-genie-v2/models"
	"github.com/ernestechie/cbt-genie-v2/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var usersColl *mongo.Collection

// Initialize mongodb client
func init()  {
	c := db.ConnectDB()

	// Set global mongo client
	client = c

	usersColl = c.Database("db").Collection("users")
}


// Return a new user with the "created_at" & "update_at" fields
func NewUser() models.User {
	return models.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Register new user
func RegisterUser (c *fiber.Ctx) error {
		var user = NewUser()

		// Parse and validate request body using utils
		if errs := utils.ParseAndValidate(c, &user); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"errors": errs,
				"message": "Error processing request",
			})
		}

		// Check if user exist with that email.
		userExistsFilter := bson.D{
			{Key: "$and",
				Value: bson.A{
					bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: strings.ToLower(user.Email)}}}},
				}},
		}

		userExistsErr := usersColl.FindOne(context.TODO(), userExistsFilter)
			if userExistsErr != nil {
				if userExistsErr.Err() == mongo.ErrNoDocuments {
					fmt.Println("No existing user found with this email.")
					} else {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"success": false,
							"message": "This email has already been used",
						})
				}
			}

		// Always convert emails to lower case to ensure consistency.
		user.Email = strings.ToLower(user.Email)

		// Create new user if the user does not exist.
		result, err := usersColl.InsertOne(context.TODO(), user)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"errors": err.Error(),
				"message": "Error creating user",
			})
		}

		userId := result.InsertedID.(bson.ObjectID)
		user.ID = userId

		// If validation passes, return the user data
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "New user created successfully!",
			"data": user,
		})
}

// Retrieve one user
func GetUser (c *fiber.Ctx) error {
	var getUser models.User
	var userId = c.Params("userId")

	userObjectId, userIdErr := bson.ObjectIDFromHex(userId);
	if userIdErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Invalid user ID, %v", userId),
		})
	}
	fmt.Println(userObjectId)
	filter := bson.M{"_id": userObjectId}
	err := usersColl.FindOne(context.TODO(), filter).Decode(&getUser)
		if err != nil {
			fmt.Println(err)

			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

	// If validation passes, return the user data
	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "User retrieved successfully!",
		"data": getUser,
	})
}

// Get all users
func GetAllUsers (c *fiber.Ctx) error {
	sort := bson.D{{Key: "createdAt", Value: 1}}
	// filter := bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 18}}}}

	cursor, err := usersColl.Find(context.TODO(), bson.M{}, options.Find().SetSort(sort))
	if err != nil {
		log.Println("usersColl.Find \n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err.Error(),
			"message": "Error getting all users",
		})
	}
	
	// Unpacks the cursor into a slice
	var results []models.User

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("cursor.All \n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors": err.Error(),
			"message": "Error getting all users",
		})
	}

	// 
	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"message": "Users retrieved successfully!",
		"data": results,
	})
}

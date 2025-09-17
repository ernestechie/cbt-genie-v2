package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UpdateBy struct {
	ID				bson.ObjectID 		`json:"id" bson:"_id,omitempty"`
	Date 			time.Time 				`json:"updated_by" bson:"updated_by"`
	UserID		bson.ObjectID 		`json:"user_id" bson:"user_id,omitempty"`
}

type User struct {
	ID           		bson.ObjectID 		`json:"id" bson:"_id,omitempty"`
	Email						string 		`json:"email" validate:"email" bson:"email"`
	FirstName				string 		`json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName				string		`json:"last_name,omitempty" bson:"last_name,omitempty"`
	Age 						int				`json:"age,omitempty" validate:"gte=12,lte=100" bson:"age,omitempty"`
	Goal						string		`json:"goal,omitempty"`
	PreferredTopics	[]string	`json:"preferred_topics,omitempty"`
	CreatedAt    time.Time		`json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time		`json:"updated_at" bson:"updated_at"`
	OtpToken    	string			`json:"otp_token" bson:"otp_token,omitempty"`
	OtpExpiry    	time.Time		`json:"otp_expiry" bson:"otp_expiry,omitempty"`
	UpdatedBy			[]UpdateBy 	`json:"updated_by,omitempty" bson:"updated_by,omitempty"`
}

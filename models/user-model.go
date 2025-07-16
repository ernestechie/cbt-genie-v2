package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           		bson.ObjectID 		`json:"id" bson:"_id,omitempty"`
	Email						string 		`json:"email" validate:"required,email" bson:"email"`
	FirstName				string 		`json:"firstName" validate:"required,min=3,max=32" bson:"firstName,omitempty"`
	LastName				string		`json:"lastName" validate:"required,min=3,max=32" bson:"lastName,omitempty"`
	Age 						int				`json:"age" validate:"required,gte=12,lte=100"`
	Goal						string		`json:"goal"`
	PreferredTopics	[]string	`json:"preferredTopics"`
	CreatedAt    time.Time		`json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time		`json:"updatedAt" bson:"updatedAt"`
}

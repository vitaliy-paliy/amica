package models

import (
  "go.mongodb.org/mongo-driver/bson/primitive" 
)

type User struct {
  ID           primitive.ObjectID `json:"_id" bson:"_id"`
  Username     string             `json:"username" bson:"username"`
  PhoneNumber  string             `json:"phone_number" bson:"phone_number"` 
  Friends      []User             `json:"friends,omitempty" bson:"friends,omitempty"`
}

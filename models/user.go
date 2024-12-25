package models

import "time"

// User represents a simple user model
type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Age       int       `json:"age" bson:"age"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}
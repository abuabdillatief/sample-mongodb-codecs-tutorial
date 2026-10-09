package models

import (
	"time"
)

type Passport struct {
    Number string `bson:"number"`
    ExpiredAt time.Time `bson:"expired_at"`
}

type User struct {
    ID          string `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Email       string             `bson:"email"`
    CreatedAt   time.Time          `bson:"created_at"`
    Passport    Passport           `bson:"passport"`
}
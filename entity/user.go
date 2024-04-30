package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct represents a user in the system with unique name, surname combination,
// unique username, and a password that should be encrypted before storage.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // Unique identifier for the user
	Name     string             `bson:"name"`          // Name of the user; should be unique with surname
	Surname  string             `bson:"surname"`       // Surname of the user; should be unique with name
	Username string             `bson:"username"`      // Username, must be unique
	Password string             `bson:"password"`      // Password, should be encrypted and 6-16 characters long
}

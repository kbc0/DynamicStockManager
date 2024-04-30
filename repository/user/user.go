package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/kbc0/DynamicStockManager/entity" // adjust the import path according to your project structure
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user entity.User) (primitive.ObjectID, error) {
	// Check if username or (name and surname) combination already exists
	exists, err := r.checkUniqueFields(user.Username, user.Name, user.Surname)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if exists {
		return primitive.NilObjectID, errors.New("username or name-surname combination already exists")
	}

	// Insert the user into the database
	result, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to convert the inserted ID")
	}

	return oid, nil
}

// checkUniqueFields ensures that the username and name-surname combination are unique
func (r *UserRepository) checkUniqueFields(username, name, surname string) (bool, error) {
	filter := bson.M{"$or": []bson.M{
		{"username": username},
		{"name": name, "surname": surname},
	}}
	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// More functions can be added here for updating, deleting, and retrieving users
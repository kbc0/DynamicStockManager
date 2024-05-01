package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kbc0/DynamicStockManager/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormRepository struct {
	collection *mongo.Collection
}

func NewFormRepository(db *mongo.Database) *FormRepository {
	return &FormRepository{
		collection: db.Collection("forms"),
	}
}

// CreateForm inserts a new form into the database, ensuring the form name is unique per user
func (r *FormRepository) CreateForm(form entity.Form) error {
	// Check for unique form name for the user
	filter := bson.M{"userId": form.UserID, "name": form.Name}
	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("form name must be unique per user")
	}

	// Insert the form
	_, err = r.collection.InsertOne(context.TODO(), form)
	return err
}

// GetFormsByUserID retrieves all forms for a specific user
func (r *FormRepository) GetFormsByUserID(userID uuid.UUID, limit int64, offset int64) ([]entity.Form, error) {
	var forms []entity.Form
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	cursor, err := r.collection.Find(context.TODO(), bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var form entity.Form
		if err := cursor.Decode(&form); err != nil {
			continue
		}
		forms = append(forms, form)
	}
	return forms, nil
}

// GetFormByID retrieves a single form by ID
func (r *FormRepository) GetFormByID(id uuid.UUID) (*entity.Form, error) {
	var form entity.Form
	if err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&form); err != nil {
		return nil, err
	}
	return &form, nil
}

// UpdateForm updates an existing form
func (r *FormRepository) UpdateForm(form entity.Form) error {
	result, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": form.ID}, bson.M{"$set": form})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no changes applied or form not found")
	}
	return nil
}

// DeleteForm deletes a form
func (r *FormRepository) DeleteForm(id uuid.UUID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

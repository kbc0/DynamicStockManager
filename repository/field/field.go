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

type FieldRepository struct {
    collection *mongo.Collection
}

func NewFieldRepository(db *mongo.Database) *FieldRepository {
    return &FieldRepository{
        collection: db.Collection("fields"),
    }
}

// CreateField inserts a new field into the database, ensuring field name uniqueness within a form
func (r *FieldRepository) CreateField(field entity.Field) error {
    // Check for unique field name within the form
    filter := bson.M{"formId": field.FormID, "name": field.Name}
    count, err := r.collection.CountDocuments(context.TODO(), filter)
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("field name must be unique within the form")
    }

    // Insert the field
    _, err = r.collection.InsertOne(context.TODO(), field)
    return err
}

// GetFieldsByFormID retrieves all fields for a specific form, sorted by the field order
func (r *FieldRepository) GetFieldsByFormID(formID uuid.UUID) ([]entity.Field, error) {
    var fields []entity.Field
    opts := options.Find().SetSort(bson.M{"order": 1}) // Sort by field order
    cursor, err := r.collection.Find(context.TODO(), bson.M{"formId": formID}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())
    for cursor.Next(context.TODO()) {
        var field entity.Field
        if err := cursor.Decode(&field); err != nil {
            continue
        }
        fields = append(fields, field)
    }
    return fields, nil
}

// GetFieldByID retrieves a single field by ID
func (r *FieldRepository) GetFieldByID(id uuid.UUID) (*entity.Field, error) {
    var field entity.Field
    if err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&field); err != nil {
        return nil, err
    }
    return &field, nil
}

// UpdateField updates an existing field
func (r *FieldRepository) UpdateField(field entity.Field) error {
    result, err := r.collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": field.ID},
        bson.M{"$set": field},
    )
    if err != nil {
        return err
    }
    if result.ModifiedCount == 0 {
        return errors.New("no changes applied or field not found")
    }
    return nil
}

// DeleteField deletes a field
func (r *FieldRepository) DeleteField(id uuid.UUID) error {
    _, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
    return err
}

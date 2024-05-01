package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/kbc0/DynamicStockManager/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StockRepository struct {
	collection *mongo.Collection
}

func NewStockRepository(db *mongo.Database) *StockRepository {
	return &StockRepository{
		collection: db.Collection("stocks"),
	}
}

func (r *StockRepository) CreateStock(stock entity.Stock) error {
	_, err := r.collection.InsertOne(context.Background(), stock)
	return err
}

func (r *StockRepository) GetStockById(id uuid.UUID) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&stock)
	if err != nil {
		return nil, err
	}
	return &stock, nil
}
func (r *StockRepository) GetAllStocksByFormId(formId uuid.UUID) ([]entity.Stock, error) {
	var stocks []entity.Stock
	filter := bson.M{"formId": formId}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var stock entity.Stock
		if err := cursor.Decode(&stock); err != nil {
			return nil, err 
		}
		stocks = append(stocks, stock)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return stocks, nil
}


func (r *StockRepository) UpdateStock(stock entity.Stock) error {
	filter := bson.M{"_id": stock.ID}
	update := bson.M{"$set": stock}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *StockRepository) DeleteStock(id uuid.UUID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

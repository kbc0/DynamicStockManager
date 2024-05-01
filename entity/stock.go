package entity

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID        uuid.UUID              `json:"id" bson:"_id"`
	FormID    uuid.UUID              `json:"formId" bson:"formId"`
	Data      map[string]interface{} `json:"data" bson:"data"` // Dynamic data storage based on form fields
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updatedAt"`
}

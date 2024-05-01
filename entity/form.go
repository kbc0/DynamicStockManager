package entity

import (
	"time"

	"github.com/google/uuid"
)

type Form struct {
    ID       uuid.UUID `json:"id" bson:"_id"`
    UserID   uuid.UUID `json:"userId" bson:"userId"`
    Name     string    `json:"name" bson:"name"`
    CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

package entity

import "github.com/google/uuid"

type FieldType string

const (
    Combobox      FieldType = "combobox"
    Text          FieldType = "text"
    Checkbox      FieldType = "checkbox"
    Number        FieldType = "number"
    NumberDecimal FieldType = "numberDecimal"
)

type Field struct {
    ID           uuid.UUID `json:"id" bson:"_id"`
    FormID       uuid.UUID `json:"formId" bson:"formId"`
    Name         string    `json:"name" bson:"name"`
    Type         FieldType `json:"type" bson:"type"`
    IsHidden     bool      `json:"isHidden" bson:"isHidden"`
    Order        int       `json:"order" bson:"order"`
    IsUnique     bool      `json:"isUnique" bson:"isUnique"`
    Options      []string  `json:"options,omitempty" bson:"options,omitempty"` // For combobox
    MinValue     *int      `json:"minValue,omitempty" bson:"minValue,omitempty"` // For number and numberDecimal
    MaxValue     *int      `json:"maxValue,omitempty" bson:"maxValue,omitempty"` // For number and numberDecimal
    DefaultValue interface{} `json:"defaultValue,omitempty" bson:"defaultValue,omitempty"` // For number, numberDecimal, and combobox
}

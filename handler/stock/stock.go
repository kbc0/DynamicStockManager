package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kbc0/DynamicStockManager/entity"
	"github.com/kbc0/DynamicStockManager/repository/stock"
	fieldRepo "github.com/kbc0/DynamicStockManager/repository/field"
)

type StockHandler struct {
	repo *repository.StockRepository
	fieldRepo *fieldRepo.FieldRepository
}

func NewStockHandler(repo *repository.StockRepository, fieldRepo *fieldRepo.FieldRepository) *StockHandler {
	return &StockHandler{
		repo: repo,
		fieldRepo: fieldRepo,
	}
}

func (h *StockHandler) AddStock(c *fiber.Ctx) error {
    formID, err := uuid.Parse(c.Params("_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form ID format"})
    }

    var data map[string]interface{}
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Retrieve fields for the form
    fields, err := h.fieldRepo.GetFieldsByFormID(formID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // Prepare to validate data and check for extraneous fields
    fieldMap := make(map[string]entity.Field)
    for _, field := range fields {
        fieldMap[field.Name] = field
    }

    // Use default values for missing fields and validate all provided fields
    for key, value := range data {
        field, exists := fieldMap[key]
        if !exists {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field provided: " + key})
        }
        if err := validateFieldData(value, field); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
    }

    // Check for required fields and use default values where appropriate
    for fieldName, field := range fieldMap {
        if _, ok := data[fieldName]; !ok {
            if !field.IsHidden { // If field is not hidden and data not provided
                if field.DefaultValue != nil {
                    data[fieldName] = field.DefaultValue // Use the default value
                } else {
                    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required field: " + fieldName})
                }
            }
        }
    }

    stock := entity.Stock{
        ID:        uuid.New(),
        FormID:    formID,
        Data:      data,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := h.repo.CreateStock(stock); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Stock added"})
}

// validateFieldData checks if the given data is valid for the field type
func validateFieldData(value interface{}, field entity.Field) error {
	switch field.Type {
	case entity.Combobox:
		// Validate that the value is one of the options and matches the default value if specified
		valStr, ok := value.(string)
		if !ok {
			return errors.New("invalid data type for combobox, expected string")
		}
		if !contains(field.Options, valStr) {
			return errors.New("value not in combobox options")
		}
	case entity.Text:
		// Validate that the value is a string and its length is within the min and max limits
		valStr, ok := value.(string)
		if !ok {
			return errors.New("invalid data type for text, expected string")
		}
		if len(valStr) < *field.MinValue || (field.MaxValue != nil && len(valStr) > *field.MaxValue) {
			return errors.New("text length outside of specified limits")
		}
	case entity.Checkbox:
		// Validate that the value is a boolean
		if _, ok := value.(bool); !ok {
			return errors.New("invalid data type for checkbox, expected boolean")
		}
	case entity.Number:
		// Validate that the value is an integer and within the specified range
		valFloat, ok := value.(float64) // JSON numbers are decoded as float64 by default
		if !ok {
			return errors.New("invalid data type for number, expected integer")
		}
		valInt := int(valFloat) // Convert float64 to int; ensure it is a natural number
		if float64(valInt) != valFloat {
			return errors.New("invalid number value, expected integer without fractional part")
		}
		if valInt < *field.MinValue || (field.MaxValue != nil && valInt > *field.MaxValue) {
			return errors.New("number outside of specified range")
		}
	case entity.NumberDecimal:
		// Validate that the value is a float and within the specified range
		valFloat, ok := value.(float64)
		if !ok {
			return errors.New("invalid data type for numberDecimal, expected decimal number")
		}
		if valFloat < float64(*field.MinValue) || (field.MaxValue != nil && valFloat > float64(*field.MaxValue)) {
			return errors.New("decimal number outside of specified range")
		}
	default:
		return errors.New("unknown field type")
	}
	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (h *StockHandler) GetStock(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("stock_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock ID format"})
	}

	stock, err := h.repo.GetStockById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Stock not found"})
	}

	return c.JSON(stock)
}
func (h *StockHandler) GetAllStocks(c *fiber.Ctx) error {
	formId, err := uuid.Parse(c.Params("_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form ID format"})
	}

	stocks, err := h.repo.GetAllStocksByFormId(formId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stocks)
}


func (h *StockHandler) UpdateStock(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("stock_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock ID format"})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	stock := entity.Stock{
		ID:        id,
		Data:      data,
		UpdatedAt: time.Now(),
	}

	if err := h.repo.UpdateStock(stock); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Stock updated"})
}

func (h *StockHandler) DeleteStock(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("stock_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock ID format"})
	}

	if err := h.repo.DeleteStock(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Stock deleted"})
}

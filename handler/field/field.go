package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kbc0/DynamicStockManager/entity"
	"github.com/kbc0/DynamicStockManager/repository/field"
)

type FieldHandler struct {
    repo *repository.FieldRepository
}

func NewFieldHandler(repo *repository.FieldRepository) *FieldHandler {
    return &FieldHandler{
        repo: repo,
    }
}

// AddFieldToForm creates a new field and adds it to a form
func (h *FieldHandler) AddFieldToForm(c *fiber.Ctx) error {
    formID, err := uuid.Parse(c.Params("_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form ID format"})
    }

    var field entity.Field
    if err := c.BodyParser(&field); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    field.FormID = formID
	field.ID= uuid.New()

    // Perform validations based on the field type
    if err := validateField(field); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := h.repo.CreateField(field); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Field added to form"})
}

// GetAllFields retrieves all fields for a form
func (h *FieldHandler) GetAllFields(c *fiber.Ctx) error {
    formID, err := uuid.Parse(c.Params("_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form ID format"})
    }

    fields, err := h.repo.GetFieldsByFormID(formID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fields)
}

// GetField retrieves a single field by ID
func (h *FieldHandler) GetField(c *fiber.Ctx) error {
    fieldID, err := uuid.Parse(c.Params("field_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field ID format"})
    }

    field, err := h.repo.GetFieldByID(fieldID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Field not found"})
    }

    return c.JSON(field)
}

// UpdateField updates a specific field
func (h *FieldHandler) UpdateField(c *fiber.Ctx) error {
    fieldID, err := uuid.Parse(c.Params("field_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field ID format"})
    }

    var field entity.Field
    if err := c.BodyParser(&field); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    field.ID = fieldID
	field.FormID,_= uuid.Parse(c.Params("_id"))
    if err := validateField(field); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := h.repo.UpdateField(field); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Field updated"})
}

// DeleteField removes a field from a form
func (h *FieldHandler) DeleteField(c *fiber.Ctx) error {
    fieldID, err := uuid.Parse(c.Params("field_id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field ID format"})
    }

    if err := h.repo.DeleteField(fieldID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Field deleted"})
}

// validateField checks for field-specific validation rules
func validateField(field entity.Field) error {
    // Add validation logic here based on field.Type
    switch field.Type {
    case entity.Combobox:
        if len(field.Options) == 0 || len(field.Options) > 10 {
            return errors.New("Combobox must have 1 to 10 options")
        }
        if field.DefaultValue == nil || !contains(field.Options, field.DefaultValue.(string)) {
            return errors.New("Default value must be one of the provided options")
        }
    case entity.Text:
        if field.MinValue != nil && field.MaxValue != nil && *field.MinValue > *field.MaxValue {
            return errors.New("Min value cannot be greater than max value")
        }
    case entity.Number, entity.NumberDecimal:
        if field.MinValue != nil && field.MaxValue != nil && *field.MinValue > *field.MaxValue {
            return errors.New("Min value cannot be greater than max value")
        }
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

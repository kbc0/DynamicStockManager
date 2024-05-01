package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kbc0/DynamicStockManager/entity"
	"github.com/kbc0/DynamicStockManager/repository/form"
	utils "github.com/kbc0/DynamicStockManager/utils"
)

type FormHandler struct {
	repo *repository.FormRepository
}

func NewFormHandler(repo *repository.FormRepository) *FormHandler {
	return &FormHandler{
		repo: repo,
	}
}

// CreateFormHandler handles the creation of a new form
func (h *FormHandler) CreateFormHandler(c *fiber.Ctx) error {
	var form entity.Form
	form.CreatedAt= time.Now()
	form.ID= uuid.New()
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Authenticate and authorize
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	form.UserID = userID

	if err := h.repo.CreateForm(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Form created"})
}

// GetFormsHandler retrieves all forms for the authenticated user
func (h *FormHandler) GetFormsHandler(c *fiber.Ctx) error {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	limit, offset := utils.ParsePagination(c)
	forms, err := h.repo.GetFormsByUserID(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(forms)
}

// GetFormHandler retrieves a single form by ID
func (h *FormHandler) GetFormHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	form, err := h.repo.GetFormByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Form not found"})
	}
	return c.JSON(form)
}

// UpdateFormHandler handles updating an existing form
func (h *FormHandler) UpdateFormHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var form entity.Form
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	form.ID = id

	if err := h.repo.UpdateForm(form); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Form updated"})
}

// DeleteFormHandler handles the deletion of a form
func (h *FormHandler) DeleteFormHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := h.repo.DeleteForm(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Form deleted"})
}

package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/kbc0/DynamicStockManager/entity"
    userRepo "github.com/kbc0/DynamicStockManager/repository/user"
    "golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    repo *userRepo.UserRepository
}

func NewUserHandler(repo *userRepo.UserRepository) *UserHandler {
    return &UserHandler{
        repo: repo,
    }
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
    var user entity.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }
    user.Password = encryptPassword(user.Password)
    id, err := h.repo.CreateUser(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id.Hex()})
}

func encryptPassword(password string) string {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return ""
    }
    return string(hashedPassword)
}

// LoginUser authenticates a user and returns a token or error
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.BodyParser(&credentials); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Fetch user from database
    user, err := h.repo.GetUserByUsername(credentials.Username)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
    }

    // Compare hashed passwords
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
    }

    // Here you would generate a token or session ID
    // For simplicity, returning a simple success message
    return c.JSON(fiber.Map{"message": "Login successful", "user": user.Username})
}


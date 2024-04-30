package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/kbc0/DynamicStockManager/entity"
    userRepo "github.com/kbc0/DynamicStockManager/repository/user"
    utils "github.com/kbc0/DynamicStockManager/utils"
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
func encryptPassword(password string) string {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return ""
    }
    return string(hashedPassword)
}

// RegisterUser function modified to include JWT token generation
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

    // Generate JWT token
    token, err := utils.GenerateToken(user.Username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id.Hex(), "token": token})
}

// LoginUser function modified to include JWT token generation
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

    // Generate JWT token
    token, err := utils.GenerateToken(user.Username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
    }

    return c.JSON(fiber.Map{"message": "Login successful", "token": token})
}

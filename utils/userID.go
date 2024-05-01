package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// ExtractUserID extracts the user ID from the JWT token in the request
func ExtractUserID(c *fiber.Ctx) (uuid.UUID, error) {
    token := c.Locals("user").(*jwt.Token)
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, errors.New("invalid token claims")
    }

    userIDStr, ok := claims["userID"].(string)
    if !ok {
        return uuid.Nil, errors.New("userID not found in token")
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, errors.New("invalid userID format")
    }

    return userID, nil
}
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// GenerateToken creates a JWT token for authenticated users including their UUID
func GenerateToken(userID uuid.UUID, username string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)

    claims["userID"] = userID.String() // Include UUID in the token
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 144).Unix() // Token expires after 144 hours

    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
// VerifyToken verifies the JWT token and returns the claims if the token is valid
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Check if the token signing method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        // Return the secret key used to sign the token
        return []byte("secret"), nil
    })
    if err != nil {
        return nil, err
    }

    // Check if the token is valid
    if !token.Valid {
        return nil, jwt.ErrTokenInvalidId
    }

    // Extract and return the claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, jwt.ErrTokenInvalidClaims
    }

    return claims, nil
}
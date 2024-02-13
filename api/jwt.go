package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tomekzakrzewski/go-movierental/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	// TODO ERROR HANDLING
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Api-Token"]
		if !ok {
			fmt.Println("token not present in the header")
			return fmt.Errorf("not authorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			fmt.Println("here")
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		// Check token expiration
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("not authorized")
		}
		// Set the current authenticated user to the context.
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	// TODO ERROR HANDLING
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("not authorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, fmt.Errorf("not authorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("not authorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("not authorized")
	}
	return claims, nil
}

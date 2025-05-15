package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:     "your-secret-key-change-in-production",
		TokenDuration: 24 * time.Hour,
	}
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, config JWTConfig) (string, error) {
	expirationTime := time.Now().Add(config.TokenDuration)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, config JWTConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func AuthMiddleware(config JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := ValidateToken(tokenString, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StudentMiddleware checks if the user has student role
func StudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "student" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Student access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CoordinatorMiddleware checks if the user has coordinator role
func CoordinatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "coordinator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Coordinator access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MentorMiddleware checks if the user has mentor role
func MentorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "mentor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Mentor access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SessionTimeoutMiddleware enforces session timeout
func SessionTimeoutMiddleware(maxSessionDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.Next()
			return
		}

		if claimsData, ok := claims.(*Claims); ok {
			issuedAt := claimsData.IssuedAt.Time
			if time.Since(issuedAt) > maxSessionDuration {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired, please login again"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

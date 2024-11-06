package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			ctx.Abort()
			return
		}

		claims, err := validateToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.UserID)
		ctx.Set("userEmail", claims.Email)
		ctx.Set("userRole", claims.Role)
		ctx.Next()
	}
}

func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// Role-based middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString("userRole")

		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

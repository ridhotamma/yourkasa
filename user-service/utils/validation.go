package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", e.Field())
			case "gt":
				errors[field] = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
			case "gte":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
			case "email":
				errors[field] = "Invalid email format"
			default:
				errors[field] = fmt.Sprintf("Invalid value for %s", e.Field())
			}
		}
	}

	return errors
}

func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if validationErrors := HandleValidationErrors(err); len(validationErrors) > 0 {
				c.JSON(400, gin.H{"errors": validationErrors})
				c.Abort()
				return
			}
		}
	}
}

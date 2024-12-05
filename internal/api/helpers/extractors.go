package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ExtractorJsonWithValidation[T any](c *gin.Context) (*T, error) {
	var result T

	if err := c.BindJSON(&result); err != nil {
		return nil, err
	}

	if err := validate.Struct(result); err != nil {
		return nil, err

	}
	return &result, nil
}

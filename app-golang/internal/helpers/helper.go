package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Error   interface{} `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}, status ...int) {
	s := 200
	if len(status) > 0 {
		s = status[0]
	}
	c.JSON(s, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, err error, status ...int) {
	s := 400
	if len(status) > 0 {
		s = status[0]
	}

	var ve validator.ValidationErrors = nil
	if errors.As(err, &ve) {
		errs := friendlyErrorMessage(ve)
		c.JSON(s, ErrorResponse{
			Success: false,
			Error:   errs,
		})
		return
	}

	c.JSON(s, ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func friendlyErrorMessage(ve validator.ValidationErrors) map[string]string {

	out := make(map[string]string)

	for _, fe := range ve {
		switch fe.Tag() {
		case "required":
			out[fe.Namespace()] = fe.Field() + " this required"
		case "min":
			out[fe.Namespace()] = fe.Field() + " must have at least " + fe.Param() + " characters"
		case "max":
			out[fe.Namespace()] = fe.Field() + " must have a maximum " + fe.Param() + " characters"
		case "len":
			out[fe.Namespace()] = fe.Field() + " must have " + fe.Param() + " characters"
		default:
			out[fe.Namespace()] = fe.Error()
		}
	}

	return out
}

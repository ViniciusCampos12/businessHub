package helpers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
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
	c.JSON(s, ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

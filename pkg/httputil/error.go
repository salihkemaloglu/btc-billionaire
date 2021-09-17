package httputil

import "github.com/gin-gonic/gin"

// NewError is return error as json response
func NewError(ctx *gin.Context, code int, err error) {
	httpError := HTTPError{
		Code:    code,
		Message: err.Error(),
	}

	ctx.JSON(code, httpError)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"unexpected error occurred"`
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response format.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeUnauthorized = 401
	CodeForbidden    = 403
)

// OK sends a success response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Fail sends an error response with the given code and message.
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: "unauthorized",
		Data:    nil,
	})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: "forbidden",
		Data:    nil,
	})
}

// NotFound sends a 404 response.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: "not found",
		Data:    nil,
	})
}

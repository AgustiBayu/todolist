package helper

import (
	"net/http"
	"todolist/internal/domain"
	"todolist/internal/exception"

	"github.com/gin-gonic/gin"
)

// harus membuat errorhandler dan diterapkan di usecase untuk menangkap error
func NewHandleError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	if e, ok := err.(*exception.ErrorHandler); ok {
		statusCode = e.Code
	}
	c.JSON(statusCode, domain.WebResponse{
		Code:    statusCode,
		Message: http.StatusText(statusCode),
		Errors:  err.Error(),
	})
}

func NewHandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, domain.WebResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	})
}

func NewHandleLoginSuccess(c *gin.Context, data interface{}, token ...string) {
	if len(token) > 0 && token[0] != "" {
		c.Header("Authorization", "Bearer "+token[0])
		c.Header("Access-Control-Expose-Headers", "Authorization") //member aksi akses
	}
	c.JSON(http.StatusOK, domain.WebResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	})
}

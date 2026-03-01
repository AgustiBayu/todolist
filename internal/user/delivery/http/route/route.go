package route

import (
	"todolist/internal/user/delivery/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, h *http.UserHandlerImpl) {
	api := r.Group("/api/v1")
	{
		api.GET("/profile", h.GetProfile)
		api.POST("/register", h.Register)
	}
}

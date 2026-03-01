package route

import (
	middleware "todolist/internal/midleware"
	"todolist/internal/user/delivery/http"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, h *http.UserHandlerImpl) {
	api := r.Group("/api/v1")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.GET("/profile", middleware.JWTMiddleware(), h.GetProfile)
	}
}

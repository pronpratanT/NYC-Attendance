package router

import (
	"hr-program/internal/user-service/handler"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, h *handler.UserHandler) *gin.Engine {
	api := r.Group("/api")
	handler.UserRoutes(api.Group("/users"), h)
	handler.AuthRoutes(api.Group("/auth"), h)

	return r
}

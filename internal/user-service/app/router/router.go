package router

import (
	"hr-program/internal/user-service/handler"

	"github.com/gin-gonic/gin"
)

func UserRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	handler.UserRoutes(api.Group("/users"), h)

	return r
}

package router

import (
	"hr-program/internal/request-service/handler"

	"github.com/gin-gonic/gin"
)

func RequestRouter(r *gin.Engine, h *handler.RequestHandler) *gin.Engine {
	api := r.Group("/api")
	handler.RequestRoutes(api.Group("/req"), h)

	return r
}

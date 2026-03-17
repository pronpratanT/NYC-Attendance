package handler

import (
	"hr-program/internal/request-service/service"
	"hr-program/shared/middleware"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Service *service.RequestService
}

func NewRequestHandler(s *service.RequestService) *RequestHandler {
	return &RequestHandler{Service: s}
}

func RequestRoutes(r *gin.RouterGroup, h *RequestHandler) {
	protected := r.Group("")
	protected.Use(middleware.JWTAuth())
	protected.GET("/ot-details/:employee_id/:date", h.GetOTDetailsByUserIDAndDate)
	protected.GET("/ot-logs/export", h.ExportOTLogsByDateRange)
}

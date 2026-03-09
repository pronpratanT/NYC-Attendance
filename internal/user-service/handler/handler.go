package handler

import (
	"hr-program/internal/user-service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func UserRoutes(r *gin.RouterGroup, h *UserHandler) {
	// r.GET("/shifts", h.GetShifts)
	r.GET("/data", h.GetData)
	// r.GET("/shift-details", h.GetShiftDetails)
	r.GET("/shift-user-date", h.GetUserShiftByUserIDAndDate)
}

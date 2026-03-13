package handler

import (
	"hr-program/internal/user-service/service"
	"hr-program/shared/middleware"

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
	r.GET("/shift-user-daterange", h.GetUserShiftByUserIDAndDateRange)
	r.GET("/shift-user-date", h.GetUserShiftByUserIDAndDate)
}

func AuthRoutes(r *gin.RouterGroup, h *UserHandler) {
	r.POST("/login", h.Login)
	r.POST("/logout", middleware.JWTAuth(), h.Logout)
}

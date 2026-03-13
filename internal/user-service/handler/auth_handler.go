package handler

import (
	"hr-program/internal/user-service/dto"
	db "hr-program/shared/connection"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	resp, err := h.Service.Login(
		req.EmployeeID,
		req.Password,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		db.ConnectRedis(),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Logout(c *gin.Context) {
	jtiValue, exists := c.Get("jti")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token context",
		})
		return
	}

	jti, ok := jtiValue.(string)
	if !ok || jti == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

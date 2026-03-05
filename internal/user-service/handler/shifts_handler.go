package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetShifts(c *gin.Context) {
	limit := 10 // You can set a default limit or get it from query parameters
	usr, err := h.Service.GetLatestShifts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve shifts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  usr,
		"total": len(usr),
	})
}

func (h *UserHandler) GetUserRaw(c *gin.Context) {
	user, err := h.Service.GetUserRaw()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  user,
		"total": len(user),
	})
}

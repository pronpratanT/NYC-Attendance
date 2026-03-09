package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// func (h *UserHandler) GetShifts(c *gin.Context) {
// 	limit := 300 // You can set a default limit or get it from query parameters
// 	usr, err := h.Service.GetLatestShifts(limit)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to retrieve shifts",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":  usr,
// 		"total": len(usr),
// 	})
// }

func (h *UserHandler) GetData(c *gin.Context) {
	data, err := h.Service.GetData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": len(data),
	})
}

// func (h *UserHandler) GetShiftDetails(c *gin.Context) {
// 	shifts, err := h.Service.GetShifts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to retrieve shift details",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":  shifts,
// 		"total": len(shifts),
// 	})
// }

func (h *UserHandler) GetUserShiftByUserIDAndDate(c *gin.Context) {
	userIDstr := c.Query("user_id")
	dateStr := c.Query("date")

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user_id",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	shifts, err := h.Service.GetUserShiftByUserIDAndDate(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user shifts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  shifts,
		"total": len(shifts),
	})
}

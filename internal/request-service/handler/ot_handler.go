package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *RequestHandler) GetOTDetailsByUserIDAndDate(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}

	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date",
		})
		return
	}

	ot, err := h.Service.AppRepo.GetOTDetailByEmployeeCodeAndDate(employeeID, date)
	if err != nil {
		// log.Printf("GetOTDetailByEmployeeCodeAndDate error: empID=%d date=%s err=%v", employeeID, date, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve OT details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  ot,
		"total": len(ot),
	})
}

func (h *RequestHandler) ExportOTLogsByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Start date and end date are required",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}

	otData, err := h.Service.ExportOTLogsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to export ot logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  otData,
		"total": len(otData),
	})
}

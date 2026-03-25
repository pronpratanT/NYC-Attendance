package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ตัวอย่าง handler แบบง่าย ถ้าต้องการดึง logs จาก AppRepo โดยตรง
func (h *AttendanceHandler) GetAttendanceLogs(c *gin.Context) {
	logs, err := h.Service.GetAttendanceLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendance logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}

func (h *AttendanceHandler) GetAttendanceDaily(c *gin.Context) {
	daily, err := h.Service.AppRepo.GetAttendanceDaily()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendance daily",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": daily,
	})
}

func (h *AttendanceHandler) GetAttendanceDailyByEmployeeID(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}

	daily, err := h.Service.AppRepo.GetAttendanceDailyByEmployeeID(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendance daily",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  daily,
		"total": len(daily),
	})
}

func (h *AttendanceHandler) GetAttendanceDailyByEmployeeIDAndDateRange(c *gin.Context) {
	employeeID, err := strconv.ParseInt(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid employee ID",
		})
		return
	}
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")

	daily, err := h.Service.AppRepo.GetAttendanceDailyByEmployeeIDAndDateRange(employeeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendance daily",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  daily,
		"total": len(daily),
	})
}

func (h *AttendanceHandler) GetAttendanceDailyByDateRange(c *gin.Context) {
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	daily, err := h.Service.AppRepo.GetAttendanceDailyByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve attendance daily",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  daily,
		"total": len(daily),
	})
}

func (h *AttendanceHandler) GetAttendanceLogsByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required in YYYY-MM-DD format",
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

	logs, err := h.Service.GetAttendanceLogsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve attendance logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": len(logs),
	})
}

func (h *AttendanceHandler) ExportAttendanceLogsTXTByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required in YYYY-MM-DD format",
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

	content, err := h.Service.ExportAttendaceLogsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to export attendance logs",
			"details": err.Error(),
		})
		return
	}

	fileName := fmt.Sprintf("attendance_logs_%s_to_%s.txt", startDate, endDate)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
}

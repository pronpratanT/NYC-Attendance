package handler

import (
	"net/http"
	"strconv"

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
		"data": daily,
	})
}

package handler

import (
	"net/http"
	"strconv"

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

	date := c.Query("date")

	ot, err := h.Service.AppRepo.GetOTDetailByEmployeeCodeAndDate(employeeID, date)
	if err != nil {
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

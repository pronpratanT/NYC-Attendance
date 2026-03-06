package handler

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

// func (h *UserHandler) GetData(c *gin.Context) {
// 	data, err := h.Service.GetData()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to retrieve data",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data":  data,
// 		"total": len(data),
// 	})
// }

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

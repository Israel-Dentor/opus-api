package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleHealth handles GET /health
func HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
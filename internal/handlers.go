package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (h *Handler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API is healthy",
	})
}
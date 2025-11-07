package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

// Healthcheck godoc
// @Summary      Health check
// @Description  Returns the API health status.
// @Tags         Health
// @ID           healthcheck
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Message: "API is healthy"})
}

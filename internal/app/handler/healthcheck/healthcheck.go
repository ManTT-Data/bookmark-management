package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	pingRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/ping"
	healthService "github.com/ManTT-Data/bookmark-management/internal/app/service/healthcheck"
)

// HealthCheck represents the interface of health check handler
type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

// healthCheckHandler represents the health check handler
type healthCheckHandler struct {
	service healthService.Healthcheck
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(serviceName, instanceID string, pingRepo pingRepo.Ping) HealthCheck {
	return &healthCheckHandler{
		service: healthService.NewHealthcheck(serviceName, instanceID, pingRepo),
	}
}

// HealthCheck   Health check function
// @Summary      Health check function
// @Description  Health check function
// @Tags         health-check
// @Accept       application/json
// @Produce      application/json
// @Success      200 {object} model.Message
// @Router       /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	msg, err := h.service.HealthCheck(c)
	if err != nil {
		log.Error().Err(err).Msg("Health-check error")
		c.JSON(http.StatusInternalServerError, struct{ Message string }{
			Message: "Instance is not ready!",
		})
		return
	}
	c.JSON(http.StatusOK, msg)
}

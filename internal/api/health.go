package api

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type HealthCheck struct{}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

// RegisterHealth register the liveness and readiness probe endpoints
func (h *HealthCheck) Register(server *echo.Echo) {
	server.GET("/liveness", h.liveness)
	server.GET("/readiness", h.readiness)
}

// liveness godoc
// @Summary Show the status of liveness server.
// @Description get the status of server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /liveness [get]
func (h *HealthCheck) liveness(c echo.Context) error {
	response := make(map[string]string)
	response["status"] = "UP"
	return c.JSON(http.StatusOK, response)
}

// readiness godoc
// @Summary Show the status readiness of server.
// @Description get the status of server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /readiness [get]
func (h *HealthCheck) readiness(c echo.Context) error {
	response := make(map[string]string)
	response["status"] = "OK"
	return c.JSON(http.StatusOK, response)
}

package controllers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check and system status endpoints
type HealthController struct {
	startTime time.Time
}

// NewHealthController creates a new HealthController instance
func NewHealthController() *HealthController {
	return &HealthController{
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
	GoVersion string `json:"goVersion"`
	NumCPU    int    `json:"numCPU"`
}

// HealthCheck godoc
// GET /api/health
// Returns server health status, uptime, and system info
func (ctrl *HealthController) HealthCheck(c *gin.Context) {
	uptime := time.Since(ctrl.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Uptime:    uptime.Round(time.Second).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		GoVersion: runtime.Version(),
		NumCPU:    runtime.NumCPU(),
	}

	c.JSON(http.StatusOK, response)
}

// Readiness godoc
// GET /api/ready
// Returns whether the server is ready to accept traffic
func (ctrl *HealthController) Readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ready":     true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

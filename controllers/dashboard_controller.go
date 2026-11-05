package controllers

import (
	"course-management/models"
	"course-management/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(dashboardService services.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

// GetDashboard godoc
// GET /api/dashboard (Admin only)
// Uses goroutines to fetch all stats concurrently
func (ctrl *DashboardController) GetDashboard(c *gin.Context) {
	start := time.Now()

	stats, err := ctrl.dashboardService.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	duration := time.Since(start)

	c.JSON(http.StatusOK, gin.H{
		"courses":      stats.TotalCourses,
		"learners":     stats.TotalLearners,
		"enrollments":  stats.TotalEnrollments,
		"revenue":      stats.TotalRevenue,
		"responseTime": duration.String(),
	})
}

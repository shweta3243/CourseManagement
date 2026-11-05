package controllers

import (
	"course-management/models"
	"course-management/repository"
	"course-management/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnrollmentController struct {
	enrollmentService services.EnrollmentService
	userRepo          repository.UserRepository
}

func NewEnrollmentController(enrollmentService services.EnrollmentService, userRepo repository.UserRepository) *EnrollmentController {
	return &EnrollmentController{
		enrollmentService: enrollmentService,
		userRepo:          userRepo,
	}
}

// Enroll godoc
// POST /api/enrollment
func (ctrl *EnrollmentController) Enroll(c *gin.Context) {
	var req models.EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// If userId not provided in body, use the authenticated user's ID
	if req.UserID == 0 {
		userID, exists := c.Get("userId")
		if exists {
			req.UserID = userID.(int)
		}
	}

	id, err := ctrl.enrollmentService.Enroll(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Enrolled successfully", "enrollmentId": id})
}

// GetMyEnrollments godoc
// GET /api/my-enrollments
func (ctrl *EnrollmentController) GetMyEnrollments(c *gin.Context) {
	userID, _ := c.Get("userId")

	enrollments, err := ctrl.enrollmentService.GetUserEnrollments(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if enrollments == nil {
		enrollments = []models.EnrollmentDetail{}
	}

	c.JSON(http.StatusOK, enrollments)
}

// GetAllEnrollments godoc
// GET /api/enrollments (Admin only)
func (ctrl *EnrollmentController) GetAllEnrollments(c *gin.Context) {
	enrollments, err := ctrl.enrollmentService.GetAllEnrollments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if enrollments == nil {
		enrollments = []models.Enrollment{}
	}

	c.JSON(http.StatusOK, enrollments)
}

// GetAllLearners godoc
// GET /api/learners (Admin only)
func (ctrl *EnrollmentController) GetAllLearners(c *gin.Context) {
	learners, err := ctrl.userRepo.FindAllLearners(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if learners == nil {
		learners = []models.User{}
	}

	// Clear passwords from response
	for i := range learners {
		learners[i].Password = ""
	}

	c.JSON(http.StatusOK, learners)
}

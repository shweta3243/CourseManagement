package controllers

import (
	"course-management/models"
	"course-management/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubjectController struct {
	subjectService services.SubjectService
}

func NewSubjectController(subjectService services.SubjectService) *SubjectController {
	return &SubjectController{subjectService: subjectService}
}

// CreateSubject godoc
// POST /api/subjects (Admin only)
func (ctrl *SubjectController) CreateSubject(c *gin.Context) {
	var req models.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := ctrl.subjectService.CreateSubject(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subject created successfully", "subjectId": id})
}

// GetAllSubjects godoc
// GET /api/subjects
func (ctrl *SubjectController) GetAllSubjects(c *gin.Context) {
	subjects, err := ctrl.subjectService.GetAllSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if subjects == nil {
		subjects = []models.Subject{}
	}

	c.JSON(http.StatusOK, subjects)
}

// DeleteSubject godoc
// DELETE /api/subjects/:id (Admin only)
func (ctrl *SubjectController) DeleteSubject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid subject ID"})
		return
	}

	if err := ctrl.subjectService.DeleteSubject(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Subject deleted successfully"})
}

// AssignSubjectToCourse godoc
// POST /api/course-subject (Admin only)
func (ctrl *SubjectController) AssignSubjectToCourse(c *gin.Context) {
	var req models.AssignSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := ctrl.subjectService.AssignToCourse(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.MessageResponse{Message: "Subject assigned to course successfully"})
}

// RemoveSubjectFromCourse godoc
// DELETE /api/course-subject/:courseId/:subjectId (Admin only)
func (ctrl *SubjectController) RemoveSubjectFromCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid course ID"})
		return
	}
	subjectID, err := strconv.Atoi(c.Param("subjectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid subject ID"})
		return
	}

	if err := ctrl.subjectService.RemoveFromCourse(c.Request.Context(), courseID, subjectID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Subject removed from course"})
}

// GetSubjectsByCourse godoc
// GET /api/courses/:id/subjects
func (ctrl *SubjectController) GetSubjectsByCourse(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid course ID"})
		return
	}

	subjects, err := ctrl.subjectService.GetSubjectsByCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if subjects == nil {
		subjects = []models.Subject{}
	}

	c.JSON(http.StatusOK, subjects)
}

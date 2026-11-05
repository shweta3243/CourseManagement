package controllers

import (
	"course-management/models"
	"course-management/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	courseService services.CourseService
}

func NewCourseController(courseService services.CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

// CreateCourse godoc
// POST /api/courses (Admin only)
func (ctrl *CourseController) CreateCourse(c *gin.Context) {
	var req models.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := ctrl.courseService.CreateCourse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully", "courseId": id})
}

// GetAllCourses godoc
// GET /api/courses
func (ctrl *CourseController) GetAllCourses(c *gin.Context) {
	// Check for search query
	keyword := c.Query("search")
	var courses []models.Course
	var err error

	if keyword != "" {
		courses, err = ctrl.courseService.SearchCourses(c.Request.Context(), keyword)
	} else {
		courses, err = ctrl.courseService.GetAllCourses(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if courses == nil {
		courses = []models.Course{}
	}

	c.JSON(http.StatusOK, courses)
}

// GetCourseByID godoc
// GET /api/courses/:id
func (ctrl *CourseController) GetCourseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid course ID"})
		return
	}

	courseWithSubjects, err := ctrl.courseService.GetCourseWithSubjects(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, courseWithSubjects)
}

// UpdateCourse godoc
// PUT /api/courses/:id (Admin only)
func (ctrl *CourseController) UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid course ID"})
		return
	}

	var req models.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := ctrl.courseService.UpdateCourse(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Course updated successfully"})
}

// DeleteCourse godoc
// DELETE /api/courses/:id (Admin only)
func (ctrl *CourseController) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid course ID"})
		return
	}

	if err := ctrl.courseService.DeleteCourse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Course deleted successfully"})
}

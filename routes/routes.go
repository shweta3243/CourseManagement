package routes

import (
	"course-management/controllers"
	"course-management/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes with proper middleware
func SetupRoutes(
	router *gin.Engine,
	authCtrl *controllers.AuthController,
	courseCtrl *controllers.CourseController,
	subjectCtrl *controllers.SubjectController,
	enrollmentCtrl *controllers.EnrollmentController,
	dashboardCtrl *controllers.DashboardController,
	healthCtrl *controllers.HealthController,
) {
	api := router.Group("/api")

	// ─── Health Check Routes (No Auth) ─────────────────────
	api.GET("/health", healthCtrl.HealthCheck)
	api.GET("/ready", healthCtrl.Readiness)

	// ─── Public Routes (No Auth) ───────────────────────────
	api.POST("/register", authCtrl.Register)
	api.POST("/login", authCtrl.Login)

	// ─── Protected Routes (JWT Required) ───────────────────
	protected := api.Group("")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// ─── Courses (Read: All users, Write: Admin only) ──
		protected.GET("/courses", courseCtrl.GetAllCourses)
		protected.GET("/courses/:id", courseCtrl.GetCourseByID)
		protected.GET("/courses/:id/subjects", subjectCtrl.GetSubjectsByCourse)

		// ─── Enrollment (Learner) ──────────────────────────
		protected.POST("/enrollment", enrollmentCtrl.Enroll)
		protected.GET("/my-enrollments", enrollmentCtrl.GetMyEnrollments)

		// ─── Admin Only Routes ─────────────────────────────
		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		{
			// Course Management
			admin.POST("/courses", courseCtrl.CreateCourse)
			admin.PUT("/courses/:id", courseCtrl.UpdateCourse)
			admin.DELETE("/courses/:id", courseCtrl.DeleteCourse)

			// Subject Management
			admin.POST("/subjects", subjectCtrl.CreateSubject)
			admin.GET("/subjects", subjectCtrl.GetAllSubjects)
			admin.DELETE("/subjects/:id", subjectCtrl.DeleteSubject)

			// Course-Subject Assignment
			admin.POST("/course-subject", subjectCtrl.AssignSubjectToCourse)
			admin.DELETE("/course-subject/:courseId/:subjectId", subjectCtrl.RemoveSubjectFromCourse)

			// Learners & Enrollments
			admin.GET("/learners", enrollmentCtrl.GetAllLearners)
			admin.GET("/enrollments", enrollmentCtrl.GetAllEnrollments)

			// Dashboard (Concurrency)
			admin.GET("/dashboard", dashboardCtrl.GetDashboard)
		}
	}
}

package main

import (
	"course-management/config"
	"course-management/controllers"
	"course-management/middleware"
	"course-management/repository"
	"course-management/routes"
	"course-management/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	config.InitDatabase()
	defer config.DB.Close()

	// ─── Repository Layer ──────────────────────────────────
	userRepo := repository.NewUserRepository(config.DB)
	courseRepo := repository.NewCourseRepository(config.DB)
	subjectRepo := repository.NewSubjectRepository(config.DB)
	enrollmentRepo := repository.NewEnrollmentRepository(config.DB)

	// ─── Service Layer (Dependency Injection) ──────────────
	authService := services.NewAuthService(userRepo)
	courseService := services.NewCourseService(courseRepo, subjectRepo)
	subjectService := services.NewSubjectService(subjectRepo, courseRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, courseRepo, userRepo)
	dashboardService := services.NewDashboardService(courseRepo, userRepo, enrollmentRepo)

	// ─── Controller Layer ──────────────────────────────────
	authCtrl := controllers.NewAuthController(authService)
	courseCtrl := controllers.NewCourseController(courseService)
	subjectCtrl := controllers.NewSubjectController(subjectService)
	enrollmentCtrl := controllers.NewEnrollmentController(enrollmentService, userRepo)
	dashboardCtrl := controllers.NewDashboardController(dashboardService)
	healthCtrl := controllers.NewHealthController()

	// ─── Gin Router Setup ──────────────────────────────────
	router := gin.Default()

	// CORS middleware for Vue frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Custom logging middleware
	router.Use(middleware.LoggingMiddleware())

	// ─── Register Routes ───────────────────────────────────
	routes.SetupRoutes(router, authCtrl, courseCtrl, subjectCtrl, enrollmentCtrl, dashboardCtrl, healthCtrl)

	// ─── Start Server ──────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

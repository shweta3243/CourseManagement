package services

import (
	"context"
	"course-management/models"
	"course-management/repository"
	"errors"
)

// EnrollmentService defines enrollment business logic
type EnrollmentService interface {
	Enroll(ctx context.Context, req *models.EnrollRequest) (int64, error)
	GetUserEnrollments(ctx context.Context, userID int) ([]models.EnrollmentDetail, error)
	GetAllEnrollments(ctx context.Context) ([]models.Enrollment, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
}

// NewEnrollmentService creates a new EnrollmentService instance
func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	userRepo repository.UserRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
	}
}

func (s *enrollmentService) Enroll(ctx context.Context, req *models.EnrollRequest) (int64, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return 0, errors.New("user not found")
	}

	// Verify course exists
	course, err := s.courseRepo.FindByID(ctx, req.CourseID)
	if err != nil {
		return 0, errors.New("course not found")
	}

	if !course.IsActive {
		return 0, errors.New("course is not active")
	}

	// Check duplicate enrollment
	exists, err := s.enrollmentRepo.Exists(ctx, req.UserID, req.CourseID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user already enrolled in this course")
	}

	enrollment := &models.Enrollment{
		UserID:   req.UserID,
		CourseID: req.CourseID,
	}

	return s.enrollmentRepo.Create(ctx, enrollment)
}

func (s *enrollmentService) GetUserEnrollments(ctx context.Context, userID int) ([]models.EnrollmentDetail, error) {
	return s.enrollmentRepo.FindByUserID(ctx, userID)
}

func (s *enrollmentService) GetAllEnrollments(ctx context.Context) ([]models.Enrollment, error) {
	return s.enrollmentRepo.FindAll(ctx)
}

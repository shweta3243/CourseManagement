package services

import (
	"context"
	"course-management/models"
	"course-management/repository"
	"errors"
)

// SubjectService defines subject business logic
type SubjectService interface {
	CreateSubject(ctx context.Context, req *models.CreateSubjectRequest) (int64, error)
	GetAllSubjects(ctx context.Context) ([]models.Subject, error)
	DeleteSubject(ctx context.Context, id int) error
	AssignToCourse(ctx context.Context, req *models.AssignSubjectRequest) error
	RemoveFromCourse(ctx context.Context, courseID, subjectID int) error
	GetSubjectsByCourse(ctx context.Context, courseID int) ([]models.Subject, error)
}

type subjectService struct {
	subjectRepo repository.SubjectRepository
	courseRepo  repository.CourseRepository
}

// NewSubjectService creates a new SubjectService instance
func NewSubjectService(subjectRepo repository.SubjectRepository, courseRepo repository.CourseRepository) SubjectService {
	return &subjectService{
		subjectRepo: subjectRepo,
		courseRepo:  courseRepo,
	}
}

func (s *subjectService) CreateSubject(ctx context.Context, req *models.CreateSubjectRequest) (int64, error) {
	if req.Name == "" {
		return 0, errors.New("subject name is required")
	}

	subject := &models.Subject{
		SubjectName: req.Name,
	}
	return s.subjectRepo.Create(ctx, subject)
}

func (s *subjectService) GetAllSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.subjectRepo.FindAll(ctx)
}

func (s *subjectService) DeleteSubject(ctx context.Context, id int) error {
	_, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("subject not found")
	}
	return s.subjectRepo.Delete(ctx, id)
}

func (s *subjectService) AssignToCourse(ctx context.Context, req *models.AssignSubjectRequest) error {
	// Verify course exists
	_, err := s.courseRepo.FindByID(ctx, req.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	// Verify subject exists
	_, err = s.subjectRepo.FindByID(ctx, req.SubjectID)
	if err != nil {
		return errors.New("subject not found")
	}

	return s.subjectRepo.AssignToCourse(ctx, req.CourseID, req.SubjectID)
}

func (s *subjectService) RemoveFromCourse(ctx context.Context, courseID, subjectID int) error {
	return s.subjectRepo.RemoveFromCourse(ctx, courseID, subjectID)
}

func (s *subjectService) GetSubjectsByCourse(ctx context.Context, courseID int) ([]models.Subject, error) {
	return s.subjectRepo.FindByCourseID(ctx, courseID)
}

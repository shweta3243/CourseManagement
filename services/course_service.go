package services

import (
	"context"
	"course-management/models"
	"course-management/repository"
	"errors"
)

// CourseService defines course business logic
type CourseService interface {
	CreateCourse(ctx context.Context, req *models.CreateCourseRequest) (int64, error)
	GetAllCourses(ctx context.Context) ([]models.Course, error)
	GetCourseByID(ctx context.Context, id int) (*models.Course, error)
	UpdateCourse(ctx context.Context, id int, req *models.UpdateCourseRequest) error
	DeleteCourse(ctx context.Context, id int) error
	SearchCourses(ctx context.Context, keyword string) ([]models.Course, error)
	GetCourseWithSubjects(ctx context.Context, id int) (*models.CourseWithSubjects, error)
}

type courseService struct {
	courseRepo  repository.CourseRepository
	subjectRepo repository.SubjectRepository
}

// NewCourseService creates a new CourseService instance
func NewCourseService(courseRepo repository.CourseRepository, subjectRepo repository.SubjectRepository) CourseService {
	return &courseService{
		courseRepo:  courseRepo,
		subjectRepo: subjectRepo,
	}
}

func (s *courseService) CreateCourse(ctx context.Context, req *models.CreateCourseRequest) (int64, error) {
	if req.Name == "" {
		return 0, errors.New("course name is required")
	}
	if req.Fee < 0 {
		return 0, errors.New("fee cannot be negative")
	}

	course := &models.Course{
		CourseName:  req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Fees:        req.Fee,
		IsActive:    true,
	}

	return s.courseRepo.Create(ctx, course)
}

func (s *courseService) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	return s.courseRepo.FindAll(ctx)
}

func (s *courseService) GetCourseByID(ctx context.Context, id int) (*models.Course, error) {
	return s.courseRepo.FindByID(ctx, id)
}

func (s *courseService) UpdateCourse(ctx context.Context, id int, req *models.UpdateCourseRequest) error {
	// Verify course exists
	_, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("course not found")
	}
	return s.courseRepo.Update(ctx, id, req)
}

func (s *courseService) DeleteCourse(ctx context.Context, id int) error {
	_, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("course not found")
	}
	return s.courseRepo.Delete(ctx, id)
}

func (s *courseService) SearchCourses(ctx context.Context, keyword string) ([]models.Course, error) {
	return s.courseRepo.Search(ctx, keyword)
}

func (s *courseService) GetCourseWithSubjects(ctx context.Context, id int) (*models.CourseWithSubjects, error) {
	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("course not found")
	}

	subjects, err := s.subjectRepo.FindByCourseID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.CourseWithSubjects{
		Course:   *course,
		Subjects: subjects,
	}, nil
}

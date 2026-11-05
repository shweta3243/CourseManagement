package repository

import (
	"context"
	"course-management/models"
	"database/sql"
)

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) (int64, error)
	FindAll(ctx context.Context) ([]models.Course, error)
	FindByID(ctx context.Context, id int) (*models.Course, error)
	Update(ctx context.Context, id int, course *models.UpdateCourseRequest) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
	Search(ctx context.Context, keyword string) ([]models.Course, error)
}

type courseRepository struct {
	db *sql.DB
}

// NewCourseRepository creates a new CourseRepository instance
func NewCourseRepository(db *sql.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *models.Course) (int64, error) {
	query := `INSERT INTO courses (course_name, description, duration, fees, is_active) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		course.CourseName, course.Description, course.Duration, course.Fees, course.IsActive,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *courseRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	query := `SELECT course_id, course_name, description, duration, fees, is_active, created_at FROM courses ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.CourseID, &c.CourseName, &c.Description, &c.Duration, &c.Fees, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *courseRepository) FindByID(ctx context.Context, id int) (*models.Course, error) {
	course := &models.Course{}
	query := `SELECT course_id, course_name, description, duration, fees, is_active, created_at FROM courses WHERE course_id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.CourseID, &course.CourseName, &course.Description, &course.Duration,
		&course.Fees, &course.IsActive, &course.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (r *courseRepository) Update(ctx context.Context, id int, req *models.UpdateCourseRequest) error {
	query := `UPDATE courses SET course_name = ?, description = ?, duration = ?, fees = ? WHERE course_id = ?`
	_, err := r.db.ExecContext(ctx, query, req.Name, req.Description, req.Duration, req.Fee, id)
	return err
}

func (r *courseRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM courses WHERE course_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *courseRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM courses`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *courseRepository) Search(ctx context.Context, keyword string) ([]models.Course, error) {
	query := `SELECT course_id, course_name, description, duration, fees, is_active, created_at 
	          FROM courses WHERE course_name LIKE ? OR description LIKE ? ORDER BY created_at DESC`
	pattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.CourseID, &c.CourseName, &c.Description, &c.Duration, &c.Fees, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

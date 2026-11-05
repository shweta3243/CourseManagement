package repository

import (
	"context"
	"course-management/models"
	"database/sql"
)

// EnrollmentRepository defines the interface for enrollment data access
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *models.Enrollment) (int64, error)
	FindByUserID(ctx context.Context, userID int) ([]models.EnrollmentDetail, error)
	FindAll(ctx context.Context) ([]models.Enrollment, error)
	Exists(ctx context.Context, userID, courseID int) (bool, error)
	Count(ctx context.Context) (int, error)
	TotalRevenue(ctx context.Context) (float64, error)
}

type enrollmentRepository struct {
	db *sql.DB
}

// NewEnrollmentRepository creates a new EnrollmentRepository instance
func NewEnrollmentRepository(db *sql.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *models.Enrollment) (int64, error) {
	query := `INSERT INTO enrollments (user_id, course_id, status) VALUES (?, ?, 'ACTIVE')`
	result, err := r.db.ExecContext(ctx, query, enrollment.UserID, enrollment.CourseID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *enrollmentRepository) FindByUserID(ctx context.Context, userID int) ([]models.EnrollmentDetail, error) {
	query := `SELECT e.enrollment_id, c.course_name, c.duration, c.fees, e.enrollment_date, e.status 
	          FROM enrollments e 
	          INNER JOIN courses c ON e.course_id = c.course_id 
	          WHERE e.user_id = ? 
	          ORDER BY e.enrollment_date DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.EnrollmentDetail
	for rows.Next() {
		var ed models.EnrollmentDetail
		if err := rows.Scan(&ed.EnrollmentID, &ed.CourseName, &ed.Duration, &ed.Fees, &ed.EnrollmentDate, &ed.Status); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, ed)
	}
	return enrollments, nil
}

func (r *enrollmentRepository) FindAll(ctx context.Context) ([]models.Enrollment, error) {
	query := `SELECT enrollment_id, user_id, course_id, enrollment_date, status FROM enrollments ORDER BY enrollment_date DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.EnrollmentID, &e.UserID, &e.CourseID, &e.EnrollmentDate, &e.Status); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, nil
}

func (r *enrollmentRepository) Exists(ctx context.Context, userID, courseID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM enrollments WHERE user_id = ? AND course_id = ?`
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(&count)
	return count > 0, err
}

func (r *enrollmentRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM enrollments`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *enrollmentRepository) TotalRevenue(ctx context.Context) (float64, error) {
	var revenue float64
	query := `SELECT COALESCE(SUM(c.fees), 0) FROM enrollments e INNER JOIN courses c ON e.course_id = c.course_id`
	err := r.db.QueryRowContext(ctx, query).Scan(&revenue)
	return revenue, err
}

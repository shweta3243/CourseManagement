package repository

import (
	"context"
	"course-management/models"
	"database/sql"
)

// SubjectRepository defines the interface for subject data access
type SubjectRepository interface {
	Create(ctx context.Context, subject *models.Subject) (int64, error)
	FindAll(ctx context.Context) ([]models.Subject, error)
	FindByID(ctx context.Context, id int) (*models.Subject, error)
	Delete(ctx context.Context, id int) error
	AssignToCourse(ctx context.Context, courseID, subjectID int) error
	RemoveFromCourse(ctx context.Context, courseID, subjectID int) error
	FindByCourseID(ctx context.Context, courseID int) ([]models.Subject, error)
}

type subjectRepository struct {
	db *sql.DB
}

// NewSubjectRepository creates a new SubjectRepository instance
func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(ctx context.Context, subject *models.Subject) (int64, error) {
	query := `INSERT INTO subjects (subject_name) VALUES (?)`
	result, err := r.db.ExecContext(ctx, query, subject.SubjectName)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *subjectRepository) FindAll(ctx context.Context) ([]models.Subject, error) {
	query := `SELECT subject_id, subject_name, created_at FROM subjects ORDER BY subject_name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.SubjectID, &s.SubjectName, &s.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

func (r *subjectRepository) FindByID(ctx context.Context, id int) (*models.Subject, error) {
	subject := &models.Subject{}
	query := `SELECT subject_id, subject_name, created_at FROM subjects WHERE subject_id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subject.SubjectID, &subject.SubjectName, &subject.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func (r *subjectRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subjects WHERE subject_id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *subjectRepository) AssignToCourse(ctx context.Context, courseID, subjectID int) error {
	query := `INSERT INTO course_subject (course_id, subject_id) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, courseID, subjectID)
	return err
}

func (r *subjectRepository) RemoveFromCourse(ctx context.Context, courseID, subjectID int) error {
	query := `DELETE FROM course_subject WHERE course_id = ? AND subject_id = ?`
	_, err := r.db.ExecContext(ctx, query, courseID, subjectID)
	return err
}

func (r *subjectRepository) FindByCourseID(ctx context.Context, courseID int) ([]models.Subject, error) {
	query := `SELECT s.subject_id, s.subject_name, s.created_at 
	          FROM subjects s 
	          INNER JOIN course_subject cs ON s.subject_id = cs.subject_id 
	          WHERE cs.course_id = ?
	          ORDER BY s.subject_name`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.SubjectID, &s.SubjectName, &s.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

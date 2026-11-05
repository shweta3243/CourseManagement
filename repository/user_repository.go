package repository

import (
	"context"
	"course-management/models"
	"database/sql"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindAllLearners(ctx context.Context) ([]models.User, error)
	CountLearners(ctx context.Context) (int, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT user_id, name, email, password, role, created_at FROM users WHERE email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT user_id, name, email, role, created_at FROM users WHERE user_id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserID, &user.Name, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindAllLearners(ctx context.Context) ([]models.User, error) {
	query := `SELECT user_id, name, email, role, created_at FROM users WHERE role = 'LEARNER' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) CountLearners(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = 'LEARNER'`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

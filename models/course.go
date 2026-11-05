package models

import "time"

type Course struct {
	CourseID    int       `json:"courseId" db:"course_id"`
	CourseName string    `json:"courseName" db:"course_name"`
	Description string   `json:"description" db:"description"`
	Duration   string    `json:"duration" db:"duration"`
	Fees       float64   `json:"fees" db:"fees"`
	IsActive   bool      `json:"isActive" db:"is_active"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type CreateCourseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Duration    string  `json:"duration" binding:"required"`
	Fee         float64 `json:"fee" binding:"required"`
}

type UpdateCourseRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	Fee         float64 `json:"fee"`
	IsActive    *bool   `json:"isActive"`
}

package models

import "time"

type Enrollment struct {
	EnrollmentID   int       `json:"enrollmentId" db:"enrollment_id"`
	UserID         int       `json:"userId" db:"user_id"`
	CourseID       int       `json:"courseId" db:"course_id"`
	EnrollmentDate time.Time `json:"enrollmentDate" db:"enrollment_date"`
	Status         string    `json:"status" db:"status"`
}

type EnrollRequest struct {
	UserID   int `json:"userId" binding:"required"`
	CourseID int `json:"courseId" binding:"required"`
}

type EnrollmentDetail struct {
	EnrollmentID   int       `json:"enrollmentId"`
	CourseName     string    `json:"courseName"`
	Duration       string    `json:"duration"`
	Fees           float64   `json:"fees"`
	EnrollmentDate time.Time `json:"enrollmentDate"`
	Status         string    `json:"status"`
}

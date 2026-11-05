package models

import "time"

type Subject struct {
	SubjectID   int       `json:"subjectId" db:"subject_id"`
	SubjectName string    `json:"subjectName" db:"subject_name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type CreateSubjectRequest struct {
	Name string `json:"name" binding:"required"`
}

type CourseSubject struct {
	ID        int `json:"id" db:"course_subject_id"`
	CourseID  int `json:"courseId" db:"course_id"`
	SubjectID int `json:"subjectId" db:"subject_id"`
}

type AssignSubjectRequest struct {
	CourseID  int `json:"courseId" binding:"required"`
	SubjectID int `json:"subjectId" binding:"required"`
}

type CourseWithSubjects struct {
	Course   Course    `json:"course"`
	Subjects []Subject `json:"subjects"`
}

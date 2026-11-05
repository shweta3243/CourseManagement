package models

type DashboardStats struct {
	TotalCourses     int     `json:"courses"`
	TotalLearners    int     `json:"learners"`
	TotalEnrollments int     `json:"enrollments"`
	TotalRevenue     float64 `json:"revenue"`
}

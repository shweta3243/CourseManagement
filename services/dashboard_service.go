package services

import (
	"context"
	"course-management/models"
	"course-management/repository"
	"fmt"
	"sync"
	"time"
)

// DashboardService defines dashboard business logic using Go concurrency
type DashboardService interface {
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
}

type dashboardService struct {
	courseRepo     repository.CourseRepository
	userRepo       repository.UserRepository
	enrollmentRepo repository.EnrollmentRepository
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(
	courseRepo repository.CourseRepository,
	userRepo repository.UserRepository,
	enrollmentRepo repository.EnrollmentRepository,
) DashboardService {
	return &dashboardService{
		courseRepo:     courseRepo,
		userRepo:       userRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

// GetDashboardStats runs 4 queries concurrently using goroutines, channels, and WaitGroup
// This reduces total response time from ~2s (sequential) to ~500ms (parallel)
func (s *dashboardService) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	// Create a context with timeout to cancel slow DB queries
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		field string
		value interface{}
		err   error
	}

	// Buffered channel to collect results from goroutines
	ch := make(chan result, 4)
	var wg sync.WaitGroup

	// Goroutine 1: Count courses
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.courseRepo.Count(ctx)
		ch <- result{field: "courses", value: count, err: err}
	}()

	// Goroutine 2: Count learners
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.userRepo.CountLearners(ctx)
		ch <- result{field: "learners", value: count, err: err}
	}()

	// Goroutine 3: Count enrollments
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.enrollmentRepo.Count(ctx)
		ch <- result{field: "enrollments", value: count, err: err}
	}()

	// Goroutine 4: Calculate total revenue
	wg.Add(1)
	go func() {
		defer wg.Done()
		revenue, err := s.enrollmentRepo.TotalRevenue(ctx)
		ch <- result{field: "revenue", value: revenue, err: err}
	}()

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from channel
	stats := &models.DashboardStats{}
	var errs []error

	for r := range ch {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", r.field, r.err))
			continue
		}

		switch r.field {
		case "courses":
			stats.TotalCourses = r.value.(int)
		case "learners":
			stats.TotalLearners = r.value.(int)
		case "enrollments":
			stats.TotalEnrollments = r.value.(int)
		case "revenue":
			stats.TotalRevenue = r.value.(float64)
		}
	}

	if len(errs) > 0 {
		return stats, fmt.Errorf("partial errors: %v", errs)
	}

	return stats, nil
}

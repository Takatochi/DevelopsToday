package services

import (
	"context"
	"sync"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
)

// SalaryUpdate represents a salary update operation
type SalaryUpdate struct {
	ID     uint    `json:"id" validate:"required"`
	Salary float64 `json:"salary" validate:"required,min=0"`
}

// BulkResult represents the result of a bulk operation
type BulkResult struct {
	Successful int      `json:"successful"`
	Failed     int      `json:"failed"`
	Errors     []string `json:"errors,omitempty"`
}

// BulkContext interface for bulk operations
type BulkContext interface {
	BulkUpdateSalary(ctx context.Context, updates []SalaryUpdate) (*BulkResult, error)
	BulkCreateCats(ctx context.Context, cats []*models.Cat) (*BulkResult, error)
}

// Bulk service implementation
type Bulk struct {
	catRepo repo.CatRepository
}

// NewBulk creates a new bulk service
func NewBulk(catRepo repo.CatRepository) BulkContext {
	return &Bulk{
		catRepo: catRepo,
	}
}

// BulkUpdateSalary updates multiple cat salaries using worker pool pattern
func (s *Bulk) BulkUpdateSalary(ctx context.Context, updates []SalaryUpdate) (*BulkResult, error) {
	if len(updates) == 0 {
		return &BulkResult{}, nil
	}

	const maxWorkers = 5
	numWorkers := maxWorkers
	if len(updates) < maxWorkers {
		numWorkers = len(updates)
	}

	// Channels for job distribution
	jobs := make(chan SalaryUpdate, len(updates))
	results := make(chan error, len(updates))

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for update := range jobs {
				err := s.catRepo.UpdateSalary(ctx, update.ID, update.Salary)
				results <- err
			}
		}()
	}

	// Send jobs to workers
	for _, update := range updates {
		jobs <- update
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	// Collect results
	var errors []string
	successful := 0
	failed := 0

	for err := range results {
		if err != nil {
			errors = append(errors, err.Error())
			failed++
		} else {
			successful++
		}
	}

	return &BulkResult{
		Successful: successful,
		Failed:     failed,
		Errors:     errors,
	}, nil
}

// BulkCreateCats creates multiple cats using worker pool pattern
func (s *Bulk) BulkCreateCats(ctx context.Context, cats []*models.Cat) (*BulkResult, error) {
	if len(cats) == 0 {
		return &BulkResult{}, nil
	}

	const maxWorkers = 5
	numWorkers := maxWorkers
	if len(cats) < maxWorkers {
		numWorkers = len(cats)
	}

	// Channels for job distribution
	jobs := make(chan *models.Cat, len(cats))
	results := make(chan error, len(cats))

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cat := range jobs {
				err := s.catRepo.Create(ctx, cat)
				results <- err
			}
		}()
	}

	// Send jobs to workers
	for _, cat := range cats {
		jobs <- cat
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	// Collect results
	var errors []string
	successful := 0
	failed := 0

	for err := range results {
		if err != nil {
			errors = append(errors, err.Error())
			failed++
		} else {
			successful++
		}
	}

	return &BulkResult{
		Successful: successful,
		Failed:     failed,
		Errors:     errors,
	}, nil
}

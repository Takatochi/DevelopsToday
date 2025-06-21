package services

import (
	"context"
	"sync"
	"time"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalCats         int     `json:"total_cats"`
	TotalMissions     int     `json:"total_missions"`
	CompletedMissions int     `json:"completed_missions"`
	TotalTargets      int     `json:"total_targets"`
	CompletedTargets  int     `json:"completed_targets"`
	AverageSalary     float64 `json:"average_salary"`
	ActiveMissions    int     `json:"active_missions"`
}

// StatsContext interface for statistics operations
type StatsContext interface {
	GetDashboard(ctx context.Context) (*DashboardStats, error)
}

// Stats service implementation
type Stats struct {
	catRepo     repo.CatRepository
	missionRepo repo.MissionRepository
}

// NewStats creates a new stats service
func NewStats(catRepo repo.CatRepository, missionRepo repo.MissionRepository) StatsContext {
	return &Stats{
		catRepo:     catRepo,
		missionRepo: missionRepo,
	}
}

// GetDashboard fetches dashboard statistics using parallel goroutines
func (s *Stats) GetDashboard(ctx context.Context) (*DashboardStats, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Result structures for goroutines
	type catResult struct {
		cats []models.Cat
		err  error
	}
	type missionResult struct {
		missions []models.Mission
		err      error
	}

	// Channels for results
	catChan := make(chan catResult, 1)
	missionChan := make(chan missionResult, 1)

	// WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch cats in parallel
	go func() {
		defer wg.Done()
		cats, err := s.catRepo.FindAll(ctx)
		catChan <- catResult{cats: cats, err: err}
	}()

	// Fetch missions in parallel
	go func() {
		defer wg.Done()
		missions, err := s.missionRepo.FindAll(ctx)
		missionChan <- missionResult{missions: missions, err: err}
	}()

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(catChan)
		close(missionChan)
	}()

	// Collect results
	var catRes catResult
	var missionRes missionResult

	// Use select to handle context cancellation
	select {
	case catRes = <-catChan:
		if catRes.err != nil {
			return nil, catRes.err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case missionRes = <-missionChan:
		if missionRes.err != nil {
			return nil, missionRes.err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Calculate statistics
	stats := &DashboardStats{
		TotalCats:     len(catRes.cats),
		TotalMissions: len(missionRes.missions),
	}

	// Calculate cat statistics
	if len(catRes.cats) > 0 {
		totalSalary := 0.0
		for _, cat := range catRes.cats {
			totalSalary += cat.Salary
		}
		stats.AverageSalary = totalSalary / float64(len(catRes.cats))
	}

	// Calculate mission and target statistics
	completedMissions := 0
	totalTargets := 0
	completedTargets := 0
	activeMissions := 0

	for _, mission := range missionRes.missions {
		if mission.Complete {
			completedMissions++
		} else if mission.CatID != nil {
			activeMissions++
		}

		totalTargets += len(mission.Targets)
		for _, target := range mission.Targets {
			if target.Complete {
				completedTargets++
			}
		}
	}

	stats.CompletedMissions = completedMissions
	stats.TotalTargets = totalTargets
	stats.CompletedTargets = completedTargets
	stats.ActiveMissions = activeMissions

	return stats, nil
}

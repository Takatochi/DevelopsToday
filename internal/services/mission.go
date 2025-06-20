package services

import (
	"context"
	"errors"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
)

type MissionContext interface {
	Create(ctx context.Context, m *models.Mission) error
	AssignCat(ctx context.Context, missionID, catID uint) error
	MarkComplete(ctx context.Context, missionID uint) error
	GetAll(ctx context.Context) ([]models.Mission, error)
	GetByID(ctx context.Context, id uint) (*models.Mission, error)
	DeleteByID(ctx context.Context, id uint) error
}

type Mission struct {
	repo repo.MissionRepository
}

func NewMission(repo repo.MissionRepository) *Mission {
	return &Mission{repo: repo}
}

func (s *Mission) Create(ctx context.Context, m *models.Mission) error {
	if len(m.Targets) < 1 || len(m.Targets) > 3 {
		return errors.New("mission must have between 1 and 3 targets")
	}
	return s.repo.Create(ctx, m)
}

func (s *Mission) AssignCat(ctx context.Context, missionID, catID uint) error {
	return s.repo.AssignCat(ctx, missionID, catID)
}

func (s *Mission) MarkComplete(ctx context.Context, missionID uint) error {
	m, err := s.repo.FindByID(ctx, missionID)
	if err != nil {
		return err
	}
	for _, t := range m.Targets {
		if !t.Complete {
			return errors.New("all targets must be completed before completing mission")
		}
	}
	return s.repo.MarkComplete(ctx, missionID)
}

func (s *Mission) GetAll(ctx context.Context) ([]models.Mission, error) {
	return s.repo.FindAll(ctx)
}

func (s *Mission) GetByID(ctx context.Context, id uint) (*models.Mission, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Mission) DeleteByID(ctx context.Context, id uint) error {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if m.CatID != nil {
		return errors.New("cannot delete assigned mission")
	}
	return s.repo.DeleteByID(ctx, id)
}

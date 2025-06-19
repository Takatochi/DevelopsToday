package services

import (
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
	"context"
	"errors"
)

type TargetContext interface {
	Add(ctx context.Context, missionID uint, t *models.Target) error
	UpdateNotes(ctx context.Context, missionID, targetID uint, notes string) error
	MarkComplete(ctx context.Context, targetID uint) error
	DeleteByID(ctx context.Context, missionID, targetID uint) error
}

type Target struct {
	targetRepo  repo.TargetRepository
	missionRepo repo.MissionRepository
}

func NewTarget(t repo.TargetRepository, m repo.MissionRepository) TargetContext {
	return &Target{
		targetRepo:  t,
		missionRepo: m,
	}
}

func (s *Target) Add(ctx context.Context, missionID uint, t *models.Target) error {
	m, err := s.missionRepo.FindByID(ctx, missionID)
	if err != nil {
		return err
	}
	if m.Complete {
		return errors.New("cannot add target to completed mission")
	}
	return s.targetRepo.AddToMission(ctx, missionID, t)
}

func (s *Target) UpdateNotes(ctx context.Context, missionID, targetID uint, notes string) error {
	m, err := s.missionRepo.FindByID(ctx, missionID)
	if err != nil {
		return err
	}
	if m.Complete {
		return errors.New("mission is completed")
	}
	for _, t := range m.Targets {
		if t.ID == targetID && t.Complete {
			return errors.New("target is completed")
		}
	}
	return s.targetRepo.UpdateNotes(ctx, targetID, notes)
}

func (s *Target) MarkComplete(ctx context.Context, targetID uint) error {
	return s.targetRepo.MarkComplete(ctx, targetID)
}

func (s *Target) DeleteByID(ctx context.Context, missionID, targetID uint) error {
	m, err := s.missionRepo.FindByID(ctx, missionID)
	if err != nil {
		return err
	}
	for _, t := range m.Targets {
		if t.ID == targetID && t.Complete {
			return errors.New("cannot delete completed target")
		}
	}
	return s.targetRepo.DeleteByID(ctx, targetID)
}

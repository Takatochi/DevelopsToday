package postgres

import (
	"context"

	"DevelopsToday/internal/models"
)

type TargetRepository struct {
	store *Repository
}

func (r *TargetRepository) AddToMission(ctx context.Context, missionID uint, target *models.Target) error {
	target.ID = missionID
	return r.store.db.WithContext(ctx).Create(target).Error
}

func (r *TargetRepository) UpdateNotes(ctx context.Context, targetID uint, notes string) error {
	return r.store.db.WithContext(ctx).
		Model(&models.Target{}).
		Where("id = ?", targetID).
		Update("notes", notes).Error
}

func (r *TargetRepository) MarkComplete(ctx context.Context, targetID uint) error {
	return r.store.db.WithContext(ctx).
		Model(&models.Target{}).
		Where("id = ?", targetID).
		Update("complete", true).Error
}

func (r *TargetRepository) DeleteByID(ctx context.Context, targetID uint) error {
	return r.store.db.WithContext(ctx).Delete(&models.Target{}, targetID).Error
}

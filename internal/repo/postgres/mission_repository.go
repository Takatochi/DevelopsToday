package postgres

import (
	"context"

	"DevelopsToday/internal/models"
)

type MissionRepository struct {
	store *Repository
}

func (r *MissionRepository) Create(ctx context.Context, mission *models.Mission) error {
	return r.store.db.WithContext(ctx).Create(mission).Error
}

func (r *MissionRepository) AssignCat(ctx context.Context, missionID, catID uint) error {
	return r.store.db.WithContext(ctx).
		Model(&models.Mission{}).
		Where("id = ?", missionID).
		Update("cat_id", catID).Error
}

func (r *MissionRepository) MarkComplete(ctx context.Context, id uint) error {
	return r.store.db.WithContext(ctx).
		Model(&models.Mission{}).
		Where("id = ?", id).
		Update("complete", true).Error
}

func (r *MissionRepository) FindByID(ctx context.Context, id uint) (*models.Mission, error) {
	var m models.Mission
	err := r.store.db.WithContext(ctx).
		Preload("Targets").
		First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MissionRepository) FindAll(ctx context.Context) ([]models.Mission, error) {
	var missions []models.Mission
	err := r.store.db.WithContext(ctx).
		Preload("Targets").
		Find(&missions).Error
	return missions, err
}

func (r *MissionRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.store.db.WithContext(ctx).Delete(&models.Mission{}, id).Error
}

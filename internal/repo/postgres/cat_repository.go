package postgres

import (
	"DevelopsToday/internal/models"
	"context"
)

type CatRepository struct {
	store *Repository
}

func (r *CatRepository) Create(ctx context.Context, cat *models.Cat) error {
	return r.store.db.WithContext(ctx).Create(cat).Error
}

func (r *CatRepository) FindAll(ctx context.Context) ([]models.Cat, error) {
	var cats []models.Cat
	err := r.store.db.WithContext(ctx).Find(&cats).Error
	return cats, err
}

func (r *CatRepository) FindByID(ctx context.Context, id uint) (*models.Cat, error) {
	var cat models.Cat
	err := r.store.db.WithContext(ctx).First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CatRepository) UpdateSalary(ctx context.Context, id uint, salary float64) error {
	return r.store.db.WithContext(ctx).Model(&models.Cat{}).
		Where("id = ?", id).
		Update("salary", salary).Error
}

func (r *CatRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.store.db.WithContext(ctx).Delete(&models.Cat{}, id).Error
}

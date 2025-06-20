package postgres

import (
	"context"

	"DevelopsToday/internal/models"

	"gorm.io/gorm"
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
	result := r.store.db.WithContext(ctx).Model(&models.Cat{}).
		Where("id = ?", id).
		Update("salary", salary)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *CatRepository) DeleteByID(ctx context.Context, id uint) error {
	result := r.store.db.WithContext(ctx).Delete(&models.Cat{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

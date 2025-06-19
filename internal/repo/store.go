package repo

import (
	"DevelopsToday/internal/models"
	"context"
)

// Repository implement from interface Store
type Store interface {
	Cat() CatRepository
	Mission() MissionRepository
	//Target() TargetRepository
	//... other entity
}

type CatRepository interface {
	Create(ctx context.Context, cat *models.Cat) error
	FindAll(ctx context.Context) ([]models.Cat, error)
	FindByID(ctx context.Context, id uint) (*models.Cat, error)
	UpdateSalary(ctx context.Context, id uint, salary float64)
	Delete(ctx context.Context, id uint) error
}
type MissionRepository interface {
	Create(ctx context.Context, mission *models.Mission) error
	FindAll(ctx context.Context) ([]models.Mission, error)
	FindByID(ctx context.Context, id uint) (*models.Mission, error)
	AssignCat(ctx context.Context, id uint, catID uint) error
	MarkComplete(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
type TargetRepository interface {
	Create(ctx context.Context, target *models.Target) error
	FindAll(ctx context.Context) ([]models.Target, error)
	FindByID(ctx context.Context, id uint) (*models.Target, error)
	UpdateNotes(ctx context.Context, id uint, notes string) error
	MarkComplete(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

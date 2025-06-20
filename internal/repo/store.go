package repo

import (
	"context"

	"DevelopsToday/internal/models"
)

// Repository implement from interface Store
type Store interface {
	Cat() CatRepository
	Mission() MissionRepository
	Target() TargetRepository
	User() UserRepository
	// ... other entity
}

type CatRepository interface {
	Create(ctx context.Context, cat *models.Cat) error
	FindAll(ctx context.Context) ([]models.Cat, error)
	FindByID(ctx context.Context, id uint) (*models.Cat, error)
	UpdateSalary(ctx context.Context, id uint, salary float64) error
	DeleteByID(ctx context.Context, id uint) error
}
type MissionRepository interface {
	Create(ctx context.Context, mission *models.Mission) error
	FindAll(ctx context.Context) ([]models.Mission, error)
	FindByID(ctx context.Context, id uint) (*models.Mission, error)
	AssignCat(ctx context.Context, id uint, catID uint) error
	MarkComplete(ctx context.Context, id uint) error
	DeleteByID(ctx context.Context, id uint) error
}
type TargetRepository interface {
	AddToMission(ctx context.Context, missionID uint, target *models.Target) error
	UpdateNotes(ctx context.Context, targetID uint, notes string) error
	MarkComplete(ctx context.Context, targetID uint) error
	DeleteByID(ctx context.Context, id uint) error
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	DeleteByID(ctx context.Context, id uint) error
	FindAll(ctx context.Context, limit, offset int) ([]*models.User, error)
}

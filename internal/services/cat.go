package services

import (
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
	"context"
)

type CatContext interface {
	Create(ctx context.Context, c *models.Cat) error
	GetAll(ctx context.Context) ([]models.Cat, error)
	GetByID(ctx context.Context, id uint) (*models.Cat, error)
	UpdateSalary(ctx context.Context, id uint, salary float64) error
	DeleteByID(ctx context.Context, id uint) error
}

type Cat struct {
	repo repo.CatRepository
}

func NewCat(repo repo.CatRepository) CatContext {
	return &Cat{repo: repo}
}

func (s *Cat) Create(ctx context.Context, c *models.Cat) error {
	return s.repo.Create(ctx, c)
}

func (s *Cat) GetAll(ctx context.Context) ([]models.Cat, error) {
	return s.repo.FindAll(ctx)
}

func (s *Cat) GetByID(ctx context.Context, id uint) (*models.Cat, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Cat) UpdateSalary(ctx context.Context, id uint, salary float64) error {
	return s.repo.UpdateSalary(ctx, id, salary)
}

func (s *Cat) DeleteByID(ctx context.Context, id uint) error {
	return s.repo.DeleteByID(ctx, id)
}

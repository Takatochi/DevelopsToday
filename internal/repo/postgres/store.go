package postgres

import (
	"DevelopsToday/internal/repo"

	"gorm.io/gorm"
)

type Repository struct {
	db                *gorm.DB
	catRepository     *CatRepository
	missionRepository *MissionRepository
	targetRepository  *TargetRepository
	userRepository    *UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Cat() repo.CatRepository {
	if r.catRepository != nil {
		return r.catRepository
	}

	r.catRepository = &CatRepository{
		store: r,
	}

	return r.catRepository
}

func (r *Repository) Mission() repo.MissionRepository {
	if r.missionRepository != nil {
		return r.missionRepository
	}

	r.missionRepository = &MissionRepository{
		store: r,
	}

	return r.missionRepository
}
func (r *Repository) Target() repo.TargetRepository {
	if r.missionRepository != nil {
		return r.targetRepository
	}

	r.targetRepository = &TargetRepository{
		store: r,
	}

	return r.targetRepository
}

func (r *Repository) User() repo.UserRepository {
	if r.userRepository != nil {
		return r.userRepository
	}

	r.userRepository = &UserRepository{
		store: r,
	}

	return r.userRepository
}

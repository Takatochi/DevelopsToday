package postgres

import (
	"context"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
)

type UserRepository struct {
	store *Repository
}

func NewUserRepository(store *Repository) repo.UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.store.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.store.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.store.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.store.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.store.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.store.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := r.store.db.WithContext(ctx).Limit(limit).Offset(offset)
	err := query.Find(&users).Error
	return users, err
}

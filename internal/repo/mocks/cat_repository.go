package mocks

import (
	"context"
	"errors"

	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

type MockCatRepository struct {
	store *Mocks
}

func (r *MockCatRepository) Create(ctx context.Context, cat *models.Cat) error {
	r.store.mutex.Lock()
	defer r.store.mutex.Unlock()

	if cat.ID == 0 {
		cat.ID = r.store.nextCatID
		r.store.nextCatID++
	}

	if _, exists := r.store.cats[cat.ID]; exists {
		return errors.New("cat with this ID already exists")
	}

	newCat := &models.Cat{
		ID:         cat.ID,
		Name:       cat.Name,
		Experience: cat.Experience,
		Breed:      cat.Breed,
		Salary:     cat.Salary,
	}

	r.store.cats[cat.ID] = newCat
	*cat = *newCat

	return nil
}

func (r *MockCatRepository) FindAll(ctx context.Context) ([]models.Cat, error) {
	r.store.mutex.RLock()
	defer r.store.mutex.RUnlock()

	cats := make([]models.Cat, 0, len(r.store.cats))
	for _, cat := range r.store.cats {
		cats = append(cats, *cat)
	}

	return cats, nil
}

func (r *MockCatRepository) FindByID(ctx context.Context, id uint) (*models.Cat, error) {
	r.store.mutex.RLock()
	defer r.store.mutex.RUnlock()

	cat, exists := r.store.cats[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	result := &models.Cat{
		ID:         cat.ID,
		Name:       cat.Name,
		Experience: cat.Experience,
		Breed:      cat.Breed,
		Salary:     cat.Salary,
	}

	return result, nil
}

func (r *MockCatRepository) UpdateSalary(ctx context.Context, id uint, salary float64) error {
	r.store.mutex.Lock()
	defer r.store.mutex.Unlock()

	cat, exists := r.store.cats[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	cat.Salary = salary
	return nil
}

func (r *MockCatRepository) DeleteByID(ctx context.Context, id uint) error {
	r.store.mutex.Lock()
	defer r.store.mutex.Unlock()

	if _, exists := r.store.cats[id]; !exists {
		return gorm.ErrRecordNotFound
	}

	delete(r.store.cats, id)
	return nil
}

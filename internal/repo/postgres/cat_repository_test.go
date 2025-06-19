package postgres

import (
	"context"
	"testing"

	"DevelopsToday/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Cat{}, &models.Mission{}, &models.Target{})
	assert.NoError(t, err)

	return db
}

func TestCatRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := &CatRepository{store: &Repository{db: db}}
	ctx := context.Background()

	t.Run("Create should create new cat", func(t *testing.T) {
		cat := &models.Cat{
			Name:       "TestCat",
			Experience: 5,
			Breed:      "Persian",
			Salary:     1000,
		}

		err := repo.Create(ctx, cat)
		assert.NoError(t, err)
		assert.NotZero(t, cat.ID)

		// Verify cat was created
		var foundCat models.Cat
		err = db.First(&foundCat, cat.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "TestCat", foundCat.Name)
	})

	t.Run("FindAll should return all cats", func(t *testing.T) {
		// Create test cats
		cats := []models.Cat{
			{Name: "Cat1", Experience: 1, Breed: "Breed1", Salary: 100},
			{Name: "Cat2", Experience: 2, Breed: "Breed2", Salary: 200},
		}
		for i := range cats {
			err := repo.Create(ctx, &cats[i])
			assert.NoError(t, err)
		}

		foundCats, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(foundCats), 2)
	})

	t.Run("FindByID should return existing cat", func(t *testing.T) {
		cat := &models.Cat{
			Name:       "FindByCat",
			Experience: 3,
			Breed:      "FindBreed",
			Salary:     300,
		}
		err := repo.Create(ctx, cat)
		assert.NoError(t, err)

		foundCat, err := repo.FindByID(ctx, cat.ID)
		assert.NoError(t, err)
		assert.Equal(t, "FindByCat", foundCat.Name)
		assert.Equal(t, cat.ID, foundCat.ID)
	})

	t.Run("FindByID should return error for non-existing cat", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateSalary should update cat salary", func(t *testing.T) {
		cat := &models.Cat{
			Name:       "SalaryCat",
			Experience: 4,
			Breed:      "SalaryBreed",
			Salary:     400,
		}
		err := repo.Create(ctx, cat)
		assert.NoError(t, err)

		err = repo.UpdateSalary(ctx, cat.ID, 500)
		assert.NoError(t, err)

		foundCat, err := repo.FindByID(ctx, cat.ID)
		assert.NoError(t, err)
		assert.Equal(t, float64(500), foundCat.Salary)
	})

	t.Run("UpdateSalary should return error for non-existing cat", func(t *testing.T) {
		err := repo.UpdateSalary(ctx, 999, 500)
		assert.Error(t, err)
	})

	t.Run("DeleteByID should delete existing cat", func(t *testing.T) {
		cat := &models.Cat{
			Name:       "DeleteCat",
			Experience: 5,
			Breed:      "DeleteBreed",
			Salary:     500,
		}
		err := repo.Create(ctx, cat)
		assert.NoError(t, err)

		err = repo.DeleteByID(ctx, cat.ID)
		assert.NoError(t, err)

		// Verify cat was deleted
		_, err = repo.FindByID(ctx, cat.ID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("DeleteByID should return error for non-existing cat", func(t *testing.T) {
		err := repo.DeleteByID(ctx, 999)
		assert.Error(t, err)
	})
}

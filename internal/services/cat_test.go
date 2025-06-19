package services

import (
	"context"
	"testing"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo/mocks"

	"gorm.io/gorm"
)

func TestCatService(t *testing.T) {
	store := mocks.NewRepository()
	catService := NewCat(store.Cat())
	ctx := context.Background()

	t.Run("Create should create new cat", func(t *testing.T) {
		cat := &models.Cat{
			Name:       "TestCat",
			Experience: 5,
			Breed:      "Persian",
			Salary:     1000,
		}

		err := catService.Create(ctx, cat)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cat.ID == 0 {
			t.Fatal("Expected cat to have ID assigned")
		}

		// Verify cat was created
		foundCat, err := catService.GetByID(ctx, cat.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if foundCat.Name != "TestCat" {
			t.Fatalf("Expected cat name 'TestCat', got '%s'", foundCat.Name)
		}
	})

	t.Run("GetAll should return all cats", func(t *testing.T) {
		cats, err := catService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(cats) < 5 { // Should have at least 5 initial cats
			t.Fatalf("Expected at least 5 cats, got %d", len(cats))
		}
	})

	t.Run("GetByID should return existing cat", func(t *testing.T) {
		cat, err := catService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cat.Name != "Whiskers" {
			t.Fatalf("Expected cat name 'Whiskers', got '%s'", cat.Name)
		}
	})

	t.Run("GetByID should return error for non-existing cat", func(t *testing.T) {
		_, err := catService.GetByID(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("UpdateSalary should update cat salary", func(t *testing.T) {
		err := catService.UpdateSalary(ctx, 1, 2000)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		cat, err := catService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cat.Salary != 2000 {
			t.Fatalf("Expected salary 2000, got %f", cat.Salary)
		}
	})

	t.Run("UpdateSalary should return error for non-existing cat", func(t *testing.T) {
		err := catService.UpdateSalary(ctx, 999, 2000)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("DeleteByID should delete existing cat", func(t *testing.T) {
		// First create a cat to delete
		cat := &models.Cat{
			Name:       "ToDelete",
			Experience: 1,
			Breed:      "Test",
			Salary:     500,
		}
		err := catService.Create(ctx, cat)
		if err != nil {
			t.Fatalf("Expected no error creating cat, got %v", err)
		}

		// Delete the cat
		err = catService.DeleteByID(ctx, cat.ID)
		if err != nil {
			t.Fatalf("Expected no error deleting cat, got %v", err)
		}

		// Verify cat was deleted
		_, err = catService.GetByID(ctx, cat.ID)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("DeleteByID should return error for non-existing cat", func(t *testing.T) {
		err := catService.DeleteByID(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}

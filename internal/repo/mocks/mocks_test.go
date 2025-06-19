package mocks

import (
	"context"
	"testing"

	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

func TestMockCatRepository(t *testing.T) {
	store := NewRepository()
	catRepo := store.Cat()
	ctx := context.Background()

	t.Run("FindAll should return initial cats", func(t *testing.T) {
		cats, err := catRepo.FindAll(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(cats) != 5 {
			t.Fatalf("Expected 5 cats, got %d", len(cats))
		}
	})

	t.Run("FindByID should return existing cat", func(t *testing.T) {
		cat, err := catRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if cat.Name != "Whiskers" {
			t.Fatalf("Expected cat name 'Whiskers', got '%s'", cat.Name)
		}
	})

	t.Run("FindByID should return error for non-existing cat", func(t *testing.T) {
		_, err := catRepo.FindByID(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("Create should add new cat", func(t *testing.T) {
		newCat := &models.Cat{
			Name:       "Fluffy",
			Experience: 3,
			Breed:      "Maine Coon",
			Salary:     1200,
		}

		err := catRepo.Create(ctx, newCat)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if newCat.ID == 0 {
			t.Fatal("Expected cat to have ID assigned")
		}

		// Verify cat was created
		foundCat, err := catRepo.FindByID(ctx, newCat.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if foundCat.Name != "Fluffy" {
			t.Fatalf("Expected cat name 'Fluffy', got '%s'", foundCat.Name)
		}
	})

	t.Run("UpdateSalary should update cat salary", func(t *testing.T) {
		err := catRepo.UpdateSalary(ctx, 1, 1500)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		cat, err := catRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if cat.Salary != 1500 {
			t.Fatalf("Expected salary 1500, got %f", cat.Salary)
		}
	})

	t.Run("DeleteByID should remove cat", func(t *testing.T) {
		err := catRepo.DeleteByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = catRepo.FindByID(ctx, 1)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}

func TestMockMissionRepository(t *testing.T) {
	store := NewRepository()
	missionRepo := store.Mission()
	ctx := context.Background()

	t.Run("FindAll should return initial missions", func(t *testing.T) {
		missions, err := missionRepo.FindAll(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(missions) != 5 {
			t.Fatalf("Expected 5 missions, got %d", len(missions))
		}
	})

	t.Run("FindByID should return mission with targets", func(t *testing.T) {
		mission, err := missionRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(mission.Targets) != 2 {
			t.Fatalf("Expected 2 targets, got %d", len(mission.Targets))
		}
	})

	t.Run("Create should add new mission with targets", func(t *testing.T) {
		newMission := &models.Mission{
			Complete: false,
			Targets: []models.Target{
				{Name: "Test Target", Country: "Ukraine", Notes: "Test notes", Complete: false},
			},
		}

		err := missionRepo.Create(ctx, newMission)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if newMission.ID == 0 {
			t.Fatal("Expected mission to have ID assigned")
		}

		// Verify mission was created with targets
		foundMission, err := missionRepo.FindByID(ctx, newMission.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(foundMission.Targets) != 1 {
			t.Fatalf("Expected 1 target, got %d", len(foundMission.Targets))
		}
	})

	t.Run("AssignCat should assign cat to mission", func(t *testing.T) {
		err := missionRepo.AssignCat(ctx, 4, 4) // Assign Felix to unassigned mission
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		mission, err := missionRepo.FindByID(ctx, 4)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if mission.CatID == nil || *mission.CatID != 4 {
			t.Fatalf("Expected cat ID 4, got %v", mission.CatID)
		}
	})

	t.Run("MarkComplete should mark mission as complete", func(t *testing.T) {
		err := missionRepo.MarkComplete(ctx, 4)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		mission, err := missionRepo.FindByID(ctx, 4)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !mission.Complete {
			t.Fatal("Expected mission to be complete")
		}
	})
}

func TestMockTargetRepository(t *testing.T) {
	store := NewRepository()
	targetRepo := store.Target()
	ctx := context.Background()

	t.Run("AddToMission should add target to existing mission", func(t *testing.T) {
		newTarget := &models.Target{
			Name:     "New Target",
			Country:  "Poland",
			Notes:    "New target notes",
			Complete: false,
		}

		err := targetRepo.AddToMission(ctx, 1, newTarget)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if newTarget.ID == 0 {
			t.Fatal("Expected target to have ID assigned")
		}

		// Verify target was added to mission
		missionRepo := store.Mission()
		mission, err := missionRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, target := range mission.Targets {
			if target.ID == newTarget.ID {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Expected target to be added to mission")
		}
	})

	t.Run("UpdateNotes should update target notes", func(t *testing.T) {
		err := targetRepo.UpdateNotes(ctx, 1, "Updated notes")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify notes were updated in mission
		missionRepo := store.Mission()
		mission, err := missionRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, target := range mission.Targets {
			if target.ID == 1 && target.Notes == "Updated notes" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Expected target notes to be updated")
		}
	})

	t.Run("MarkComplete should mark target as complete", func(t *testing.T) {
		err := targetRepo.MarkComplete(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify target was marked complete in mission
		missionRepo := store.Mission()
		mission, err := missionRepo.FindByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, target := range mission.Targets {
			if target.ID == 1 && target.Complete {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Expected target to be marked complete")
		}
	})
}

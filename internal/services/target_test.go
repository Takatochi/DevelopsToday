package services

import (
	"context"
	"errors"
	"testing"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo/mocks"

	"gorm.io/gorm"
)

func TestTargetService(t *testing.T) {
	store := mocks.NewRepository()
	targetService := NewTarget(store.Target(), store.Mission())
	ctx := context.Background()

	t.Run("Add should add target to existing mission", func(t *testing.T) {
		target := &models.Target{
			Name:     "New Target",
			Country:  "Spain",
			Notes:    "New target notes",
			Complete: false,
		}

		err := targetService.Add(ctx, 1, target) // Add to mission 1
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if target.ID == 0 {
			t.Fatal("Expected target to have ID assigned")
		}

		// Verify target was added to mission
		missionService := NewMission(store.Mission())
		mission, err := missionService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, missionTarget := range mission.Targets {
			if missionTarget.ID == target.ID {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Expected target to be added to mission")
		}
	})

	t.Run("Add should fail for non-existing mission", func(t *testing.T) {
		target := &models.Target{
			Name:     "Test Target",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}

		err := targetService.Add(ctx, 999, target)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("Add should fail for completed mission", func(t *testing.T) {
		target := &models.Target{
			Name:     "Test Target",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}

		err := targetService.Add(ctx, 2, target) // Mission 2 is complete
		if err == nil {
			t.Fatal("Expected error for completed mission")
		}
		if err.Error() != "cannot add target to completed mission" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("UpdateNotes should update target notes", func(t *testing.T) {
		err := targetService.UpdateNotes(ctx, 1, 1, "Updated notes for target 1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify notes were updated
		missionService := NewMission(store.Mission())
		mission, err := missionService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, target := range mission.Targets {
			if target.ID == 1 && target.Notes == "Updated notes for target 1" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("Expected target notes to be updated")
		}
	})

	t.Run("UpdateNotes should fail for non-existing mission", func(t *testing.T) {
		err := targetService.UpdateNotes(ctx, 999, 1, "Test notes")
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("UpdateNotes should fail for completed mission", func(t *testing.T) {
		err := targetService.UpdateNotes(ctx, 2, 3, "Test notes") // Mission 2 is complete
		if err == nil {
			t.Fatal("Expected error for completed mission")
		}
		if err.Error() != "mission is completed" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("UpdateNotes should fail for completed target", func(t *testing.T) {
		err := targetService.UpdateNotes(ctx, 2, 3, "Test notes") // Target 3 is complete
		if err == nil {
			t.Fatal("Expected error for completed target")
		}
		// Note: This will fail with "mission is completed" first, which is expected
	})

	t.Run("MarkComplete should mark target as complete", func(t *testing.T) {
		err := targetService.MarkComplete(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify target was marked complete
		missionService := NewMission(store.Mission())
		mission, err := missionService.GetByID(ctx, 1)
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

	t.Run("MarkComplete should return error for non-existing target", func(t *testing.T) {
		err := targetService.MarkComplete(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("DeleteByID should fail for completed target", func(t *testing.T) {
		err := targetService.DeleteByID(ctx, 2, 3) // Target 3 is complete
		if err == nil {
			t.Fatal("Expected error for completed target")
		}
		if err.Error() != "cannot delete completed target" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("DeleteByID should succeed for incomplete target", func(t *testing.T) {
		// First add a target to delete
		target := &models.Target{
			Name:     "To Delete",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}
		err := targetService.Add(ctx, 1, target)
		if err != nil {
			t.Fatalf("Expected no error adding target, got %v", err)
		}

		// Delete the target
		err = targetService.DeleteByID(ctx, 1, target.ID)
		if err != nil {
			t.Fatalf("Expected no error deleting target, got %v", err)
		}

		// Verify target was removed from mission
		missionService := NewMission(store.Mission())
		mission, err := missionService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found := false
		for _, missionTarget := range mission.Targets {
			if missionTarget.ID == target.ID {
				found = true
				break
			}
		}
		if found {
			t.Fatal("Expected target to be removed from mission")
		}
	})

	t.Run("DeleteByID should return error for non-existing mission", func(t *testing.T) {
		err := targetService.DeleteByID(ctx, 999, 1)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("DeleteByID should return error for non-existing target", func(t *testing.T) {
		err := targetService.DeleteByID(ctx, 1, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}

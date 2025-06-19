package services

import (
	"context"
	"testing"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo/mocks"

	"gorm.io/gorm"
)

func TestMissionService(t *testing.T) {
	store := mocks.NewRepository()
	missionService := NewMission(store.Mission())
	ctx := context.Background()

	t.Run("Create should create mission with valid targets", func(t *testing.T) {
		mission := &models.Mission{
			Complete: false,
			Targets: []models.Target{
				{Name: "Test Target 1", Country: "Ukraine", Notes: "Test notes 1", Complete: false},
				{Name: "Test Target 2", Country: "Poland", Notes: "Test notes 2", Complete: false},
			},
		}

		err := missionService.Create(ctx, mission)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if mission.ID == 0 {
			t.Fatal("Expected mission to have ID assigned")
		}

		// Verify mission was created
		foundMission, err := missionService.GetByID(ctx, mission.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(foundMission.Targets) != 2 {
			t.Fatalf("Expected 2 targets, got %d", len(foundMission.Targets))
		}
	})

	t.Run("Create should fail with no targets", func(t *testing.T) {
		mission := &models.Mission{
			Complete: false,
			Targets:  []models.Target{},
		}

		err := missionService.Create(ctx, mission)
		if err == nil {
			t.Fatal("Expected error for mission with no targets")
		}
		if err.Error() != "mission must have between 1 and 3 targets" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("Create should fail with more than 3 targets", func(t *testing.T) {
		mission := &models.Mission{
			Complete: false,
			Targets: []models.Target{
				{Name: "Target 1", Country: "Country 1", Notes: "Notes 1", Complete: false},
				{Name: "Target 2", Country: "Country 2", Notes: "Notes 2", Complete: false},
				{Name: "Target 3", Country: "Country 3", Notes: "Notes 3", Complete: false},
				{Name: "Target 4", Country: "Country 4", Notes: "Notes 4", Complete: false},
			},
		}

		err := missionService.Create(ctx, mission)
		if err == nil {
			t.Fatal("Expected error for mission with more than 3 targets")
		}
		if err.Error() != "mission must have between 1 and 3 targets" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("GetAll should return all missions", func(t *testing.T) {
		missions, err := missionService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(missions) < 5 { // Should have at least 5 initial missions
			t.Fatalf("Expected at least 5 missions, got %d", len(missions))
		}
	})

	t.Run("GetByID should return existing mission", func(t *testing.T) {
		mission, err := missionService.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(mission.Targets) != 2 {
			t.Fatalf("Expected 2 targets, got %d", len(mission.Targets))
		}
	})

	t.Run("GetByID should return error for non-existing mission", func(t *testing.T) {
		_, err := missionService.GetByID(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("AssignCat should assign cat to mission", func(t *testing.T) {
		err := missionService.AssignCat(ctx, 4, 4) // Assign Felix to unassigned mission
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		mission, err := missionService.GetByID(ctx, 4)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if mission.CatID == nil || *mission.CatID != 4 {
			t.Fatalf("Expected cat ID 4, got %v", mission.CatID)
		}
	})

	t.Run("AssignCat should return error for non-existing mission", func(t *testing.T) {
		err := missionService.AssignCat(ctx, 999, 1)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("AssignCat should return error for non-existing cat", func(t *testing.T) {
		err := missionService.AssignCat(ctx, 4, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("MarkComplete should fail if not all targets are complete", func(t *testing.T) {
		err := missionService.MarkComplete(ctx, 1) // Mission 1 has incomplete targets
		if err == nil {
			t.Fatal("Expected error for mission with incomplete targets")
		}
		if err.Error() != "all targets must be completed before completing mission" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("MarkComplete should succeed if all targets are complete", func(t *testing.T) {
		err := missionService.MarkComplete(ctx, 2) // Mission 2 has all targets complete
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		mission, err := missionService.GetByID(ctx, 2)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !mission.Complete {
			t.Fatal("Expected mission to be complete")
		}
	})

	t.Run("DeleteByID should fail for assigned mission", func(t *testing.T) {
		err := missionService.DeleteByID(ctx, 1) // Mission 1 is assigned to a cat
		if err == nil {
			t.Fatal("Expected error for assigned mission")
		}
		if err.Error() != "cannot delete assigned mission" {
			t.Fatalf("Expected specific error message, got %v", err)
		}
	})

	t.Run("DeleteByID should succeed for unassigned mission", func(t *testing.T) {
		// First create an unassigned mission
		mission := &models.Mission{
			Complete: false,
			CatID:    nil,
			Targets: []models.Target{
				{Name: "To Delete", Country: "Test", Notes: "Test", Complete: false},
			},
		}
		err := missionService.Create(ctx, mission)
		if err != nil {
			t.Fatalf("Expected no error creating mission, got %v", err)
		}

		// Delete the mission
		err = missionService.DeleteByID(ctx, mission.ID)
		if err != nil {
			t.Fatalf("Expected no error deleting mission, got %v", err)
		}

		// Verify mission was deleted
		_, err = missionService.GetByID(ctx, mission.ID)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("DeleteByID should return error for non-existing mission", func(t *testing.T) {
		err := missionService.DeleteByID(ctx, 999)
		if err != gorm.ErrRecordNotFound {
			t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}

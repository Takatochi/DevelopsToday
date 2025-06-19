package mocks

import (
	"context"

	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

type MockMissionRepository struct {
	store *Mocks
}

func (m *MockMissionRepository) Create(ctx context.Context, mission *models.Mission) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	if mission.ID == 0 {
		mission.ID = m.store.nextMissionID
		m.store.nextMissionID++
	}

	newMission := &models.Mission{
		ID:       mission.ID,
		CatID:    mission.CatID,
		Complete: mission.Complete,
		Targets:  make([]models.Target, len(mission.Targets)),
	}

	for i, target := range mission.Targets {
		if target.ID == 0 {
			target.ID = m.store.nextTargetID
			m.store.nextTargetID++
		}
		target.MissionID = mission.ID

		newMission.Targets[i] = target

		targetCopy := &models.Target{
			ID:        target.ID,
			Name:      target.Name,
			Country:   target.Country,
			Notes:     target.Notes,
			Complete:  target.Complete,
			MissionID: target.MissionID,
		}
		m.store.targets[target.ID] = targetCopy
	}

	m.store.missions[mission.ID] = newMission
	*mission = *newMission

	return nil
}

func (m *MockMissionRepository) FindAll(ctx context.Context) ([]models.Mission, error) {
	m.store.mutex.RLock()
	defer m.store.mutex.RUnlock()

	missions := make([]models.Mission, 0, len(m.store.missions))
	for _, mission := range m.store.missions {
		missionCopy := m.copyMissionWithTargets(mission)
		missions = append(missions, *missionCopy)
	}

	return missions, nil
}

func (m *MockMissionRepository) FindByID(ctx context.Context, id uint) (*models.Mission, error) {
	m.store.mutex.RLock()
	defer m.store.mutex.RUnlock()

	mission, exists := m.store.missions[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return m.copyMissionWithTargets(mission), nil
}

func (m *MockMissionRepository) AssignCat(ctx context.Context, missionID, catID uint) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	mission, exists := m.store.missions[missionID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Перевіряємо, чи існує кіт
	if _, catExists := m.store.cats[catID]; !catExists {
		return gorm.ErrRecordNotFound
	}

	mission.CatID = &catID
	return nil
}

func (m *MockMissionRepository) MarkComplete(ctx context.Context, id uint) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	mission, exists := m.store.missions[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	mission.Complete = true
	return nil
}

func (m *MockMissionRepository) DeleteByID(ctx context.Context, id uint) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	mission, exists := m.store.missions[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Видаляємо всі цілі цієї місії
	for _, target := range mission.Targets {
		delete(m.store.targets, target.ID)
	}

	delete(m.store.missions, id)
	return nil
}

// copyMissionWithTargets створює повну копію місії з усіма цілями
func (m *MockMissionRepository) copyMissionWithTargets(mission *models.Mission) *models.Mission {
	result := &models.Mission{
		ID:       mission.ID,
		CatID:    mission.CatID,
		Complete: mission.Complete,
		Targets:  make([]models.Target, 0),
	}

	// Знаходимо всі цілі для цієї місії
	for _, target := range m.store.targets {
		if target.MissionID == mission.ID {
			targetCopy := models.Target{
				ID:        target.ID,
				Name:      target.Name,
				Country:   target.Country,
				Notes:     target.Notes,
				Complete:  target.Complete,
				MissionID: target.MissionID,
			}
			result.Targets = append(result.Targets, targetCopy)
		}
	}

	return result
}

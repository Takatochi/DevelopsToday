package mocks

import (
	"context"

	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

type MockTargetRepository struct {
	store *Mocks
}

func (m *MockTargetRepository) AddToMission(ctx context.Context, missionID uint, target *models.Target) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	// Перевіряємо, чи існує місія
	mission, exists := m.store.missions[missionID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Присвоюємо ID цілі, якщо його немає
	if target.ID == 0 {
		target.ID = m.store.nextTargetID
		m.store.nextTargetID++
	}

	target.MissionID = missionID

	// Створюємо копію цілі
	newTarget := &models.Target{
		ID:        target.ID,
		Name:      target.Name,
		Country:   target.Country,
		Notes:     target.Notes,
		Complete:  target.Complete,
		MissionID: target.MissionID,
	}

	// Зберігаємо ціль
	m.store.targets[target.ID] = newTarget

	// Додаємо ціль до місії
	mission.Targets = append(mission.Targets, *newTarget)

	*target = *newTarget // Оновлюємо оригінальний об'єкт

	return nil
}

func (m *MockTargetRepository) UpdateNotes(ctx context.Context, targetID uint, notes string) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	target, exists := m.store.targets[targetID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	target.Notes = notes

	// Також оновлюємо в місії
	if mission, missionExists := m.store.missions[target.MissionID]; missionExists {
		for i, missionTarget := range mission.Targets {
			if missionTarget.ID == targetID {
				mission.Targets[i].Notes = notes
				break
			}
		}
	}

	return nil
}

func (m *MockTargetRepository) MarkComplete(ctx context.Context, targetID uint) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	target, exists := m.store.targets[targetID]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	target.Complete = true

	// Також оновлюємо в місії
	if mission, missionExists := m.store.missions[target.MissionID]; missionExists {
		for i, missionTarget := range mission.Targets {
			if missionTarget.ID == targetID {
				mission.Targets[i].Complete = true
				break
			}
		}
	}

	return nil
}

func (m *MockTargetRepository) DeleteByID(ctx context.Context, id uint) error {
	m.store.mutex.Lock()
	defer m.store.mutex.Unlock()

	target, exists := m.store.targets[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Видаляємо з місії
	if mission, missionExists := m.store.missions[target.MissionID]; missionExists {
		for i, missionTarget := range mission.Targets {
			if missionTarget.ID == id {
				// Видаляємо елемент зі слайсу
				mission.Targets = append(mission.Targets[:i], mission.Targets[i+1:]...)
				break
			}
		}
	}

	// Видаляємо з загального сховища
	delete(m.store.targets, id)

	return nil
}

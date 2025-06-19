package mocks

import (
	"sync"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
)

type Mocks struct {
	cats                  map[uint]*models.Cat
	missions              map[uint]*models.Mission
	targets               map[uint]*models.Target
	mockCatRepository     *MockCatRepository
	mockMissionRepository *MockMissionRepository
	mockTargetRepository  *MockTargetRepository
	mutex                 sync.RWMutex
	nextCatID             uint
	nextMissionID         uint
	nextTargetID          uint
}

func NewRepository() *Mocks {
	m := &Mocks{
		cats:          make(map[uint]*models.Cat),
		missions:      make(map[uint]*models.Mission),
		targets:       make(map[uint]*models.Target),
		nextCatID:     1,
		nextMissionID: 1,
		nextTargetID:  1,
	}

	// Додаємо початкові тестові дані
	m.seedData()

	return m
}

// seedData додає початкові тестові дані
func (m *Mocks) seedData() {
	// Додаємо котів-шпигунів
	cats := []*models.Cat{
		{ID: 1, Name: "Whiskers", Experience: 5, Breed: "Bengal", Salary: 1000},
		{ID: 2, Name: "Shadow", Experience: 2, Breed: "Siamese", Salary: 800},
		{ID: 3, Name: "Mittens", Experience: 8, Breed: "Persian", Salary: 1500},
		{ID: 4, Name: "Felix", Experience: 3, Breed: "Maine Coon", Salary: 900},
		{ID: 5, Name: "Luna", Experience: 6, Breed: "Russian Blue", Salary: 1200},
	}

	for _, cat := range cats {
		m.cats[cat.ID] = cat
	}
	m.nextCatID = 6

	// Додаємо цілі
	targets := []*models.Target{
		{ID: 1, Name: "Mr. Brie", Country: "France", Notes: "Cheese thefts in Paris", Complete: false, MissionID: 1},
		{ID: 2, Name: "Dr. Dre", Country: "Germany", Notes: "Suspicious barking in Berlin", Complete: false, MissionID: 1},
		{ID: 3, Name: "Agent Smith", Country: "USA", Notes: "Matrix activities completed", Complete: true, MissionID: 2},
		{ID: 4, Name: "The Fisherman", Country: "Japan", Notes: "Illegal fishing operations", Complete: false, MissionID: 3},
		{ID: 5, Name: "Sushi Master", Country: "Japan", Notes: "Suspicious sushi activities", Complete: true, MissionID: 3},
		{ID: 6, Name: "Ninja Cat", Country: "Japan", Notes: "Stealth training required", Complete: false, MissionID: 3},
		{ID: 7, Name: "The Yarn Ball", Country: "Canada", Notes: "Missing yarn investigation", Complete: false, MissionID: 4},
		{ID: 8, Name: "Laser Pointer", Country: "UK", Notes: "Mysterious red dot sightings", Complete: false, MissionID: 5},
		{ID: 9, Name: "Cardboard Box", Country: "UK", Notes: "Suspicious packaging", Complete: false, MissionID: 5},
	}

	for _, target := range targets {
		m.targets[target.ID] = target
	}
	m.nextTargetID = 10

	// Додаємо місії
	catID1 := uint(1) // Whiskers
	catID2 := uint(2) // Shadow
	catID3 := uint(3) // Mittens
	catID5 := uint(5) // Luna

	missions := []*models.Mission{
		{
			ID:       1,
			CatID:    &catID1,
			Complete: false,
			Targets:  []models.Target{*targets[0], *targets[1]}, // Mr. Brie, Dr. Dre
		},
		{
			ID:       2,
			CatID:    &catID2,
			Complete: true,
			Targets:  []models.Target{*targets[2]}, // Agent Smith
		},
		{
			ID:       3,
			CatID:    &catID3,
			Complete: false,
			Targets:  []models.Target{*targets[3], *targets[4], *targets[5]}, // Japan targets
		},
		{
			ID:       4,
			CatID:    nil, // Unassigned
			Complete: false,
			Targets:  []models.Target{*targets[6]}, // The Yarn Ball
		},
		{
			ID:       5,
			CatID:    &catID5,
			Complete: false,
			Targets:  []models.Target{*targets[7], *targets[8]}, // UK targets
		},
	}

	for _, mission := range missions {
		m.missions[mission.ID] = mission
	}
	m.nextMissionID = 6
}

// Допоміжні методи для тестування
func (m *Mocks) AddCat(cat *models.Cat) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if cat.ID == 0 {
		cat.ID = m.nextCatID
		m.nextCatID++
	}
	m.cats[cat.ID] = cat
}

func (m *Mocks) AddMission(mission *models.Mission) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if mission.ID == 0 {
		mission.ID = m.nextMissionID
		m.nextMissionID++
	}
	m.missions[mission.ID] = mission
}

func (m *Mocks) AddTarget(target *models.Target) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if target.ID == 0 {
		target.ID = m.nextTargetID
		m.nextTargetID++
	}
	m.targets[target.ID] = target
}
func (m *Mocks) Cat() repo.CatRepository {
	if m.mockCatRepository != nil {
		return m.mockCatRepository
	}

	m.mockCatRepository = &MockCatRepository{
		store: m,
	}

	return m.mockCatRepository
}

func (m *Mocks) Mission() repo.MissionRepository {
	if m.mockMissionRepository != nil {
		return m.mockMissionRepository
	}

	m.mockMissionRepository = &MockMissionRepository{
		store: m,
	}

	return m.mockMissionRepository
}
func (m *Mocks) Target() repo.TargetRepository {
	if m.mockTargetRepository != nil {
		return m.mockTargetRepository
	}

	m.mockTargetRepository = &MockTargetRepository{
		store: m,
	}

	return m.mockTargetRepository
}

package repo

import (
	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// Створюємо користувачів (тільки якщо їх ще немає)
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{Username: "admin", Email: "admin@spycats.com", Password: "admin123", Role: "admin"},
			{Username: "agent", Email: "agent@spycats.com", Password: "agent123", Role: "user"},
			{Username: "manager", Email: "manager@spycats.com", Password: "manager123", Role: "manager"},
		}
		for _, user := range users {
			var existingUser models.User
			if err := db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err = db.Create(&user).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
	}

	// Створюємо котів-шпигунів
	cats := []models.Cat{
		{Name: "Whiskers", Experience: 5, Breed: "Bengal", Salary: 1000},
		{Name: "Shadow", Experience: 2, Breed: "Siamese", Salary: 800},
		{Name: "Mittens", Experience: 8, Breed: "Persian", Salary: 1500},
		{Name: "Felix", Experience: 3, Breed: "Maine Coon", Salary: 900},
		{Name: "Luna", Experience: 6, Breed: "Russian Blue", Salary: 1200},
	}
	if err := db.Create(&cats).Error; err != nil {
		return err
	}

	// Створюємо місії з цілями
	missions := []models.Mission{
		{
			CatID:    &cats[0].ID, // Whiskers
			Complete: false,
			Targets: []models.Target{
				{Name: "Mr. Brie", Country: "France", Notes: "Cheese thefts in Paris", Complete: false},
				{Name: "Dr. Dre", Country: "Germany", Notes: "Suspicious barking in Berlin", Complete: false},
			},
		},
		{
			CatID:    &cats[1].ID, // Shadow
			Complete: true,
			Targets: []models.Target{
				{Name: "Agent Smith", Country: "USA", Notes: "Matrix activities completed", Complete: true},
			},
		},
		{
			CatID:    &cats[2].ID, // Mittens
			Complete: false,
			Targets: []models.Target{
				{Name: "The Fisherman", Country: "Japan", Notes: "Illegal fishing operations", Complete: false},
				{Name: "Sushi Master", Country: "Japan", Notes: "Suspicious sushi activities", Complete: true},
				{Name: "Ninja Cat", Country: "Japan", Notes: "Stealth training required", Complete: false},
			},
		},
		{
			CatID:    nil, // Unassigned mission
			Complete: false,
			Targets: []models.Target{
				{Name: "The Yarn Ball", Country: "Canada", Notes: "Missing yarn investigation", Complete: false},
			},
		},
		{
			CatID:    &cats[4].ID, // Luna
			Complete: false,
			Targets: []models.Target{
				{Name: "Laser Pointer", Country: "UK", Notes: "Mysterious red dot sightings", Complete: false},
				{Name: "Cardboard Box", Country: "UK", Notes: "Suspicious packaging activities", Complete: false},
			},
		},
	}

	for _, mission := range missions {
		if err := db.Create(&mission).Error; err != nil {
			return err
		}
	}

	return nil
}

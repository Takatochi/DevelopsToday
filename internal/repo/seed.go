package repo

import (
	"DevelopsToday/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	cats := []models.Cat{
		{Name: "Whiskers", Experience: 5, Breed: "beng", Salary: 1000},
		{Name: "Shadow", Experience: 2, Breed: "sibe", Salary: 800},
	}
	if err := db.Create(&cats).Error; err != nil {
		return err
	}

	mission := models.Mission{
		CatID:    &cats[0].ID,
		Complete: false,
		Targets: []models.Target{
			{Name: "Mr. Mouse", Country: "France", Notes: "Cheese thefts", Complete: false},
			{Name: "Dr. Dog", Country: "Germany", Notes: "Suspicious barking", Complete: false},
		},
	}
	if err := db.Create(&mission).Error; err != nil {
		return err
	}

	return nil
}

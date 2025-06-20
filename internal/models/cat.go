package models

// Cat represents a cat entity
// @Description Cat entity
type Cat struct {
	ID         uint    `gorm:"primaryKey" json:"id" example:"1"`
	Name       string  `json:"name" validate:"required,min=1,max=100" example:"Whiskers"`
	Breed      string  `json:"breed" validate:"required,min=1,max=50" example:"Bengal"`
	Experience int     `json:"experience" validate:"min=0,max=50" example:"5"`
	Salary     float64 `json:"salary" validate:"min=0" example:"1000.50"`
}

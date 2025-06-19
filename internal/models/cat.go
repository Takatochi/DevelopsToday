package models

// Cat represents a cat entity
// @Description Cat entity
type Cat struct {
	Name       string  `json:"name"`
	Breed      string  `json:"breed"`
	ID         uint    `gorm:"primaryKey" json:"id" example:"1"`
	Experience int     `json:"experience"`
	Salary     float64 `json:"salary"`
}

package models

// Cat represents a cat entity
// @Description Cat entity
type Cat struct {
	ID         uint    `gorm:"primaryKey" json:"id" example:"1"`
	Name       string  `json:"name"`
	Experience int     `json:"experience"`
	Breed      string  `json:"breed"`
	Salary     float64 `json:"salary"`
}

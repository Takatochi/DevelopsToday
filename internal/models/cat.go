package models

// Cat represents a cat entity
// @Description Cat entity
type Cat struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Experience int     `json:"experience"`
	Breed      string  `json:"breed"`
	Salary     float64 `json:"salary"`
}

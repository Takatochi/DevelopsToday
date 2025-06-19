package models

// Mission represents a mission entity with targets
// @Description Mission entity with assigned targets and cat
type Mission struct {
	ID       uint     `gorm:"primaryKey" json:"id" example:"1"`
	Targets  []Target `json:"targets"`
	CatID    *uint    `json:"cat_id,omitempty" example:"1"`
	Complete bool     `json:"complete" example:"false"`
}

package models

// Mission represents a mission entity with targets
// @Description Mission entity with assigned targets and cat
type Mission struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	CatID    *uint    `json:"cat_id"`
	Complete bool     `json:"complete"`
	Targets  []Target `gorm:"foreignKey:MissionID;constraint:OnDelete:CASCADE;" json:"targets"`
}

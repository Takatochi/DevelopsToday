package models

// Mission represents a mission entity with targets
// @Description Mission entity with assigned targets and cat
type Mission struct {
	CatID    *uint    `json:"cat_id"`
	Targets  []Target `gorm:"foreignKey:MissionID;constraint:OnDelete:CASCADE;" json:"targets"`
	ID       uint     `gorm:"primaryKey" json:"id"`
	Complete bool     `json:"complete"`
}

package models

// Target represents a target entity within a mission
// @Description Target entity that needs to be completed
type Target struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Complete  bool   `json:"complete"`
	MissionID uint   `json:"-"`
}

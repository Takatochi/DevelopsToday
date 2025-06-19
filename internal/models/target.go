package models

// Target represents a target entity within a mission
// @Description Target entity that needs to be completed
type Target struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	ID        uint   `gorm:"primaryKey" json:"id"`
	MissionID uint   `json:"-"`
	Complete  bool   `json:"complete"`
}

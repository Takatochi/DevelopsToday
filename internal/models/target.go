package models

// Target represents a target entity within a mission
// @Description Target entity that needs to be completed
type Target struct {
	ID        uint   `gorm:"primaryKey" json:"id" example:"1"`
	Name      string `json:"name" validate:"required,min=1,max=100" example:"John Doe"`
	Country   string `json:"country" validate:"required,min=1,max=50" example:"Ukraine"`
	Notes     string `json:"notes" validate:"required,min=1,max=500" example:"Target usually visits gym at 6 PM"`
	MissionID uint   `json:"-"`
	Complete  bool   `json:"complete" example:"false"`
}

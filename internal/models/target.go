package models

// Target represents a target entity within a mission
// @Description Target entity that needs to be completed
type Target struct {
	ID       uint   `json:"id" example:"1"`
	Name     string `json:"name" example:"Mr. Smith"`
	Notes    string `json:"notes,omitempty" example:"Target usually visits gym at 6 PM"`
	Complete bool   `json:"complete" example:"false"`
}

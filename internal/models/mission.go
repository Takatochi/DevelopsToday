package models

type Mission struct {
	ID       uint     `json:"id"`
	Targets  []Target `json:"targets"`
	CatID    *uint    `json:"cat_id,omitempty"`
	Complete bool     `json:"complete"`
}

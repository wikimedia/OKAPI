package models

import "time"

type timestamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (model *timestamp) OnUpdate() {
	model.UpdatedAt = time.Now()
}

func (model *timestamp) OnInsert() {
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
}

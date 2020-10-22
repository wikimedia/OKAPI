package models

import "time"

type timestampModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (timestamp *timestampModel) OnUpdate() {
	timestamp.UpdatedAt = time.Now()
}

func (timestamp *timestampModel) OnInsert() {
	timestamp.CreatedAt = time.Now()
	timestamp.UpdatedAt = time.Now()
}

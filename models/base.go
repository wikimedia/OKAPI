package models

import (
	"time"
)

type baseModel struct {
	ID        int       `pg:",pk" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (base *baseModel) IsUpdate() bool {
	return base.ID > 0
}

func (base *baseModel) OnUpdate() {
	base.UpdatedAt = time.Now()
}

func (base *baseModel) OnInsert() {
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
}

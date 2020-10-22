package models

type baseModel struct {
	timestampModel
	ID int `pg:",pk" json:"id"`
}

func (base *baseModel) IsUpdate() bool {
	return base.ID > 0
}

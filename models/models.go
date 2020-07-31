package models

// Model interface for all models
type Model interface {
	IsUpdate() bool
}

// Save function to save model data
func Save(model Model) error {
	if model.IsUpdate() {
		return DB().Update(model)
	}

	return DB().Insert(model)
}

// Delete function to delete the model
func Delete(model Model) error {
	return DB().Delete(model)
}

package models

// Model interface for all models
type Model interface {
	IsUpdate() bool
}

// Save function to save model data
func Save(model Model) error {
	if model.IsUpdate() {
		_, err := DB().Model(model).WherePK().Update()
		return err
	}

	_, err := DB().Model(model).Insert()
	return err
}

// Delete function to delete the model
func Delete(model Model) error {
	_, err := DB().Model(model).WherePK().Delete()
	return err
}

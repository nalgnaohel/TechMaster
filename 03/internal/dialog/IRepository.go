package dialog

import "03/internal/models"

type DialogRepository interface {
	GetByID(ID int64) (*models.Dialog, error)
	GetAll() (*models.DialogList, error)
	Create(dialog *models.Dialog) (*models.Dialog, error)
}

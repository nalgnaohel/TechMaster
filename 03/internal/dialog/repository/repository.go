package repository

import (
	"03/internal/dialog"
	"03/internal/models"

	"gorm.io/gorm"
)

type dialogRepository struct {
	db *gorm.DB
}

func NewDialogRepository(db *gorm.DB) dialog.DialogRepository {
	return &dialogRepository{
		db: db,
	}
}

func (dr *dialogRepository) GetByID(ID int64) (*models.Dialog, error) {
	var dialog models.Dialog

	err := dr.db.Where("id = ?", ID).First(&dialog).Error
	if err != nil {
		return nil, err
	}
	return &dialog, nil
}

func (dr *dialogRepository) GetAll() (*models.DialogList, error) {
	var dialogs models.DialogList

	err := dr.db.Find(&dialogs).Error
	if err != nil {
		return nil, err
	}
	return &dialogs, nil
}

func (dr *dialogRepository) Create(dialog *models.Dialog) (*models.Dialog, error) {
	err := dr.db.Table("dialog").Create(dialog).Error
	if err != nil {
		return nil, err
	}
	return dialog, nil
}

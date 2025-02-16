package repository

import (
	"03/internal/models"
	"03/internal/word"
	"errors"

	"gorm.io/gorm"
)

type wordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) word.WordRepository {
	return &wordRepository{
		db: db,
	}
}

func (wr *wordRepository) GetByID(ID int64) (*models.Word, error) {
	var word models.Word

	err := wr.db.Where("id = ?", ID).First(&word).Error
	if err != nil {
		return nil, err
	}
	return &word, nil
}

func (wr *wordRepository) GetByContent(Content string) (*models.WordList, error) {
	var words models.WordList

	err := wr.db.Where("content = ?", Content).Find(&words).Error
	if err != nil {
		return nil, err
	}
	return &words, nil
}

func (wr *wordRepository) GetAll() (*models.WordList, error) {
	var words models.WordList

	err := wr.db.Find(&words).Error
	if err != nil {
		return nil, err
	}
	return &words, nil
}

func (wr *wordRepository) Create(word *models.Word) (*models.Word, error) {
	err := wr.db.Table("word").Where("content = ? AND translate = ?", word.Content, word.Translate).First(&word).Error
	if err == nil {
		return word, errors.New("word already exists")
	}
	er := wr.db.Table("word").Create(word).Error
	if er != nil {
		return nil, er
	}
	return word, nil
}

func (wr *wordRepository) AddDialogWord(dialogID int64, wordID int64) error {
	err := wr.db.Exec("INSERT INTO word_dialog (dialog_id, word_id) VALUES (?, ?)", dialogID, wordID).Error
	if err != nil {
		return err
	}
	return nil
}

package word

import "03/internal/models"

type WordRepository interface {
	// GetWord returns a word from the repository
	GetByID(ID int64) (*models.Word, error)
	GetByContent(Content string) (*models.WordList, error)
	GetAll() (*models.WordList, error)
	Create(word *models.Word) (*models.Word, error)
	AddDialogWord(dialogID int64, wordID int64) error
}

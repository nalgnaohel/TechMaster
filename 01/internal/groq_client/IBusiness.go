package groq_client

import (
	"app01/internal/models"
)

type GroqBusiness interface {
	ChatCompletion(groqClient *models.GroqClient, prompt string) (*string, error)
}

package groq_client

import (
	"03/internal/models"
)

type GroqBusiness interface {
	ChatCompletion(groqClient *models.GroqClient, prompt string) (*string, error)
}

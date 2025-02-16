package groq_client

import (
	"03/internal/models"
	"context"
)

type GroqBusiness interface {
	ChatCompletion(c context.Context, groqClient *models.GroqClient, prompt string) (*string, *string, error)
}

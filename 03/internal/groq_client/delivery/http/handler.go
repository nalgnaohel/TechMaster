package http

import (
	"03/internal/groq_client"
	"03/internal/models"
	"log"
	"os"

	"github.com/kataras/iris/v12"
)

type groqHandler struct {
	GroqBusiness groq_client.GroqBusiness
}

func NewGroqHandler(ip iris.Party, groqBusiness groq_client.GroqBusiness) {
	handler := &groqHandler{
		GroqBusiness: groqBusiness,
	}
	ip.Post("/submit", handler.ChatCompletion)
}

func (gh *groqHandler) ChatCompletion(ctx iris.Context) {
	prompt := ctx.FormValue("prompt")
	log.Println("Prompt: ", prompt)

	apiKey := os.Getenv("GROQ_API_KEY")
	groqClient := &models.GroqClient{
		ApiKey: apiKey,
	}

	result, err := gh.GroqBusiness.ChatCompletion(groqClient, prompt)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("FF2: " + err.Error())
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.View("index.html", result)
}

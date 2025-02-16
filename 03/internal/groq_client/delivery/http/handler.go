package http

import (
	"03/internal/groq_client"
	"03/internal/models"
	"context"
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

	c := ctx.Request().Context()
	if c == nil {
		c = context.Background()
	}

	jsonWords, dialog, err := gh.GroqBusiness.ChatCompletion(c, groqClient, prompt)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("FF2: " + err.Error())
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.ViewData("Prompt", prompt)
	ctx.ViewData("JsonWords", *jsonWords)
	ctx.ViewData("Dialog", *dialog)
	ctx.View("index.html")
}

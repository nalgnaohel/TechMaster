package main

import (
	"app01/internal/models"
	"fmt"
	"log"
	"os"

	groqBusiness "app01/internal/groq_client/business"

	iris "github.com/kataras/iris/v12"
)

func main() {
	//Initialize Groq API key
	groqApiKey := os.Getenv("GROQ_API_KEY")
	fmt.Println("Groq API key: ", groqApiKey)
	// Initialize Iris application
	app := iris.New()

	// Set the views directory
	app.RegisterView(iris.HTML("./views", ".html"))

	// Register a route to serve the HTML form
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	// Register a route to handle form submission
	app.Post("/submit", func(ctx iris.Context) {
		prompt := ctx.FormValue("prompt")
		log.Println("Prompt: ", prompt)

		apiKey := os.Getenv("GROQ_API_KEY")
		groqClient := &models.GroqClient{
			ApiKey: apiKey,
		}

		gb := groqBusiness.NewGroqBusiness()
		result, err := gb.ChatCompletion(groqClient, prompt)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString("FF2: " + err.Error())
			return
		}
		resultText := ""
		if result != nil {
			resultText = *result
		}

		// Render the result in the HTML template
		ctx.ViewData("Result", resultText)
		ctx.ViewData("Prompt", prompt)
		ctx.View("index.html")
	})

	app.Listen(":8081")
}

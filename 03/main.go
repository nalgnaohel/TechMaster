package main

import (
	"fmt"
	"os"

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
		question := ctx.FormValue("question")
		ctx.View("index.html")
		ctx.HTML("Received question: " + question)
	})

	app.Listen(":8080")
}

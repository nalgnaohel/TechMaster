package main

import (
	"encoding/json"
	iris "github.com/kataras/iris/v12"
	"fmt"
)

func main() {
	//Initialize Iris application
	app := iris.New()

	//Register a route
	app.Get("/hello", func(ctx iris.Context) {
		//Create a map
		data := map[string]string{"message": "Hello, World!"}

		//Marshal the map into JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		//Set the response header
		ctx.Header("Content-Type", "application/json")

		//Write the JSON data to the response
		ctx.Write(jsonData)
	})

	app.Listen(":8080")

}
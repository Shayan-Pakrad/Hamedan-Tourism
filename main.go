package main

import "example.com/go-htmx-tailwindcss/application"

func main() {
	app := application.New()
	app.Setup()
	app.Start()
}

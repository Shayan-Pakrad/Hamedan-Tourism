package main

import "example.com/hamedan-tourism/application"

func main() {
	app := application.New()
	app.Setup()
	app.Start()
}

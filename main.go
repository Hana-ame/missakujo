package main

import "github.com/Hana-ame/missakujo/backend"

func main() {
	app := backend.App()

	// app.Listen(":3000")
	app.Listen("127.23.0.1:8080")
}

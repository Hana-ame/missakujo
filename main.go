package main

import "github.com/Hana-ame/missakujo/backend"

func main() {
	app := backend.App()

	app.Listen(":3000")
}

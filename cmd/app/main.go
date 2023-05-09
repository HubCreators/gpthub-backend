package main

import "auth/internal/app"

const configDir = "configs"

func main() {
	app.Run(configDir)
}

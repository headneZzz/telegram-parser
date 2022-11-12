// Binary authentic implements authentication example for user using terminal.
package main

import (
	"telegram-parser/internal/app"
	"telegram-parser/internal/config"
)

func main() {
	config.Config()
	app.Run()
}

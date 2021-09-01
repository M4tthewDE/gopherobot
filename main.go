package main

import (
	"de.com.fdm/gopherobot/bot"
	"de.com.fdm/gopherobot/config"
)

func main() {
	config := config.GetConfig()

	bot := bot.NewBot(config)
	bot.Run()
}

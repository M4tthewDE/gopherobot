package main

import (
	"log"

	"de.com.fdm/gopherobot/bot"
	"de.com.fdm/gopherobot/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := config.GetConfig()

	bot := bot.NewBot(config)
	bot.Run()
}

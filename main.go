package main

import (
	"log"

	"de.com.fdm/gopherobot/bot"
	"de.com.fdm/gopherobot/config"
	"github.com/davidbyttow/govips/v2/vips"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	vips.Startup(nil)
	defer vips.Shutdown()

	config := config.GetConfig()

	bot := bot.NewBot(config)
	bot.Run()
}

package main

import (
	"time"

	"de.com.fdm/gopherobot/bot"
	"de.com.fdm/gopherobot/config"
	"github.com/gempir/go-twitch-irc/v2"
)

var client *twitch.Client
var StartTime time.Time
var Channels []string
var Config config.Config

func main() {
	config := config.GetConfig()

	bot := bot.NewBot(config)
	bot.Run()
}

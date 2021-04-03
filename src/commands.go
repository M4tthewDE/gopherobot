package main

import (
	"github.com/gempir/go-twitch-irc/v2"
)

func EchoCommand(message twitch.PrivateMessage, client *twitch.Client) {
	payload := message.Message[5:]
	client.Say(message.Channel, payload)
}

package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"os"
	"strings"
)

var client *twitch.Client

func main() {
	client = twitch.NewClient("gopherobot", "oauth:"+os.Getenv("TWITCH_TOKEN"))

	client.OnPrivateMessage(onMessage)

	client.Join("gopherobot", "matthewde", "turtoise")

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if prefix == ";" {
		doCommand(message)
	}
}

func doCommand(message twitch.PrivateMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		EchoCommand(message, client)
	}
}

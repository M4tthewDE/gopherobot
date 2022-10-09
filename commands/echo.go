package commands

import "github.com/gempir/go-twitch-irc/v2"

func Echo(message twitch.PrivateMessage) string {
	if len(message.Message) < 6 {
		return ""
	}

	return message.Message[6:]
}

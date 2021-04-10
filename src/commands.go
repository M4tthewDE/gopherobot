package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"strconv"
	"strings"
	"time"
)

func EchoCommand(message twitch.PrivateMessage, client *twitch.Client) {
	payload := message.Message[6:]
	client.Say(message.Channel, payload)
}

func UserIdCommand(message twitch.PrivateMessage, client *twitch.Client) {
	args := strings.Split(message.Message[4:], " ")
	id, err := GetUserID(args[0])
	if err != nil {
		client.Say(message.Channel, `Couldn't find User-ID for "`+args[0]+`"`)
	} else {
		client.Say(message.Channel, "User-ID of "+args[0]+" is "+strconv.Itoa(id))
	}
}

func UserCommand(message twitch.PrivateMessage, client *twitch.Client) {
	args := strings.Split(message.Message[6:], " ")
	id, err := strconv.Atoi(args[0])
	if err != nil {
		client.Say(message.Channel, "Invalid User-ID "+args[0])
	}

	user, err := GetUser(id)
	if err != nil {
		client.Say(message.Channel, `Couldn't find User for "`+args[0]+`"`)
	} else {
		client.Say(message.Channel, "Username for "+args[0]+" is "+user)
	}
}

func AddFollowAlertCommand(message twitch.PrivateMessage, client *twitch.Client, host string) {
	args := strings.Split(message.Message[16:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		client.Say(message.Channel, `Couldn't find User-ID for "`+args[0]+`"`)
	} else {
		err = RegisterWebhook(id, host, message.Channel, message.User.Name)
		if err != nil {
			client.Say(message.Channel, "Error adding follow alert!")
			return
		}
		client.Say(message.Channel, "Added follow alert for "+args[0]+"!")
	}
}

func RemoveFollowAlertCommand(message twitch.PrivateMessage, client *twitch.Client, host string) {
	args := strings.Split(message.Message[19:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		client.Say(message.Channel, `Couldn't find User-ID for "`+args[0]+`"`)
	} else {
		err = RemoveWebhook(id, host)
		if err != nil {
			client.Say(message.Channel, "Error removing follow alert!")
		}
		client.Say(message.Channel, "Removed follow alert for "+args[0]+"!")
	}
}

func PingCommand(message twitch.PrivateMessage, client *twitch.Client, startTime time.Time) {
	payload := "Pong! Uptime: " + time.Since(startTime).String() + "!"
	client.Say(message.Channel, payload)
}

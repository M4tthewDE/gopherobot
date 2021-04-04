package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"strconv"
	"strings"
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

func AddFollowAlertCommand(message twitch.PrivateMessage, client *twitch.Client) {
	args := strings.Split(message.Message[16:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		client.Say(message.Channel, `Couldn't find User-ID for "`+args[0]+`"`)
	} else {
		RegisterWebhook(id)
	}
}

func RemoveFollowAlertCommand(message twitch.PrivateMessage, client *twitch.Client) {
	args := strings.Split(message.Message[19:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		client.Say(message.Channel, `Couldn't find User-ID for "`+args[0]+`"`)
	} else {
		DeleteWebhook(id)
	}
}

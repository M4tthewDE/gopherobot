package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"strconv"
	"strings"
	"time"
)

func EchoCommand(message twitch.PrivateMessage) string {
	return message.Message[6:]
}

func UserIdCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[4:], " ")
	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		return "User-ID of " + args[0] + " is " + strconv.Itoa(id)
	}
}

func UserCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[6:], " ")
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "Invalid User-ID " + args[0]
	}

	user, err := GetUser(id)
	if err != nil {
		return `Couldn't find User for "` + args[0] + `"`
	} else {
		return "Username for " + args[0] + " is " + user
	}
}

func AddFollowAlertCommand(message twitch.PrivateMessage, host string) string {
	args := strings.Split(message.Message[16:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		err = RegisterWebhook(id, host, message.Channel, message.User.Name)
		if err != nil {
			return "Error adding follow alert!"
		}
		return "Added follow alert for " + args[0] + "!"
	}
}

func RemoveFollowAlertCommand(message twitch.PrivateMessage, host string) string {
	args := strings.Split(message.Message[19:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		err = RemoveWebhook(id, host)
		if err != nil {
			return "Error removing follow alert!"
		}
		return "Removed follow alert for " + args[0] + "!"
	}
}

func PingCommand(message twitch.PrivateMessage, startTime time.Time) string {
	return "Pong! Uptime: " + time.Since(startTime).String() + "!"
}

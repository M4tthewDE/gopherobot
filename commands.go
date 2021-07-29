package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/hako/durafmt"
)

func EchoCommand(message twitch.PrivateMessage) string {
	return message.Message[6:]
}

func UserIdCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 4 {
		return `No user provided`
	}
	args := strings.Split(message.Message[4:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		return "User-ID of " + args[0] + " is " + id
	}
}

func UserCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 6 {
		return `No ID provided`
	}
	args := strings.Split(message.Message[6:], " ")

	id := args[0]

	user, err := GetUser(id)
	if err != nil {
		return `Couldn't find User for "` + args[0] + `"`
	} else {
		return "Username for " + args[0] + " is " + user
	}
}

func AddFollowAlertCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[16:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		err = RegisterWebhook(id, message.Channel, message.User.Name)
		if err != nil {
			return "Error adding follow alert!"
		}
		return "Added follow alert for " + args[0] + "!"
	}
}

func RemoveFollowAlertCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[19:], " ")

	id, err := GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	} else {
		err = RemoveWebhook(id, message.User.Name, message.Channel)
		if err != nil {
			return "Error removing follow alert!"
		}
		return "Removed follow alert for " + args[0] + "!"
	}
}

func GetFollowAlertsCommand() string {
	followWebhook, err := GetWebhooks()
	if err != nil {
		return "Error getting Followalerts"
	}
	payload := "Total alerts: " + strconv.Itoa(len(followWebhook.Data)) + " Channels: "

	var users []string
	for _, webhook := range followWebhook.Data {
		id := webhook.Condition.BroadcasterUserID
		user, _ := GetUser(id)
		users = append(users, user)
	}

	for _, user := range users {
		payload = payload + user + " "
	}

	return payload
}

func PingCommand(message twitch.PrivateMessage) string {
	uptime := time.Since(StartTime)
	result := "Pong! Uptime: " + durafmt.Parse(uptime).LimitFirstN(2).String() + ","

	api_uptime, err := GetApiUptime()
	if err != nil {
		result = result + " API-Uptime: Unavailable monkaS"
	} else {
		result = result + " API-Uptime: " + api_uptime + ","
	}

	result = result + " Commit: " + Conf.Git.Commit + ","
	result = result + " Branch: " + Conf.Git.Branch + ","

	latency := GetLatency(message)
	result = result + " Latency to tmi: " + latency

	return result
}

func RawMsgCommand(raw_message string) string {
	return UploadToHaste(raw_message)
}

func TmpJoinCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 9 {
		return `No channel provided`
	}
	args := strings.Split(message.Message[9:], " ")
	channel := args[0]

	client.Join(channel)
	Channels = append(Channels, channel)
	return "Joined #" + channel
}

func TmpLeaveCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 10 {
		return `No channel provided`
	}
	args := strings.Split(message.Message[10:], " ")
	channel := args[0]

	client.Depart(channel)

	var index int
	for i, c := range Channels {
		if c == channel {
			index = i
		}
	}
	Channels = append(Channels[:index], Channels[index+1:]...)

	return "Left #" + channel
}

func GetChannelsCommand(message twitch.PrivateMessage) string {
	if message.Channel != "matthewde" && message.Channel != "gopherobot" {
		return "Command not available in this channel to prevent pings"
	}

	var result string
	for _, channel := range Channels {
		result = result + ", " + channel
	}
	result = "Joined Channels: [" + result[2:] + "]"
	return result
}

func UrlEncodeCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 11 {
		return "Nothing to encode"
	}
	content := message.Message[11:]
	return url.QueryEscape(content)
}

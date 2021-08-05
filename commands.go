package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/hako/durafmt"
)

type CommandHandler struct {
	botStartTime time.Time
	branch       string
	commit       string
}

func (cmdHandler *CommandHandler) EchoCommand(message twitch.PrivateMessage) string {
	return message.Message[6:]
}

func (cmdHandler *CommandHandler) UserIdCommand(message twitch.PrivateMessage) string {
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

func (cmdHandler *CommandHandler) UserCommand(message twitch.PrivateMessage) string {
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

func (cmdHandler *CommandHandler) AddFollowAlertCommand(message twitch.PrivateMessage) string {
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

func (cmdHandler *CommandHandler) RemoveFollowAlertCommand(message twitch.PrivateMessage) string {
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

func (cmdHandler *CommandHandler) GetFollowAlertsCommand() string {
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

func (cmdHandler *CommandHandler) PingCommand(message twitch.PrivateMessage) string {
	uptime := time.Since(cmdHandler.botStartTime)
	result := "Pong! Uptime: " + durafmt.Parse(uptime).LimitFirstN(2).String() + ","

	api_uptime, err := GetApiUptime()
	if err != nil {
		result = result + " API-Uptime: Unavailable monkaS"
	} else {
		result = result + " API-Uptime: " + api_uptime + ","
	}

	result = result + " Commit: " + cmdHandler.commit + ","
	result = result + " Branch: " + cmdHandler.branch + ","

	latency := GetLatency(message)
	result = result + " Latency to tmi: " + latency

	return result
}

func (cmdHandler *CommandHandler) RawMsgCommand(raw_message string) string {
	return UploadToHaste(raw_message)
}

func (cmdHandler *CommandHandler) TmpJoinCommand(message twitch.PrivateMessage, bot_channels *[]string, bot_client *twitch.Client) string {
	if len(message.Message) < 9 {
		return `No channel provided`
	}
	args := strings.Split(message.Message[9:], " ")
	channel := args[0]

	bot_client.Join(channel)
	*bot_channels = append(*bot_channels, channel)
	return "Joined #" + channel
}

func (cmdHandler *CommandHandler) TmpLeaveCommand(message twitch.PrivateMessage, bot_channels *[]string, bot_client *twitch.Client) string {
	if len(message.Message) < 10 {
		return `No channel provided`
	}
	args := strings.Split(message.Message[10:], " ")
	channel := args[0]

	bot_client.Depart(channel)

	var index int
	for i, c := range *bot_channels {
		if c == channel {
			index = i
		}
	}
	*bot_channels = append((*bot_channels)[:index], (*bot_channels)[index+1:]...)

	return "Left #" + channel
}

func (cmdHandler *CommandHandler) GetChannelsCommand(message twitch.PrivateMessage, bot_channels *[]string) string {
	if message.Channel != "matthewde" && message.Channel != "gopherobot" {
		return "Command not available in this channel to prevent pings"
	}

	var result string
	for _, channel := range *bot_channels {
		result = result + ", " + channel
	}
	result = "Joined Channels: [" + result[2:] + "]"
	return result
}

func (cmdHandler *CommandHandler) UrlEncodeCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 11 {
		return "Nothing to encode"
	}
	content := message.Message[11:]
	return url.QueryEscape(content)
}

func (cmdHandler *CommandHandler) UrlDecodeCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 11 {
		return "Nothing to decode"
	}
	content := message.Message[11:]

	result, err := url.QueryUnescape(content)
	if err != nil {
		log.Println(err)
		return "Error decoding"
	}
	return result
}

func (cmdHandler *CommandHandler) StreamInfoCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 12 {
		return "No channel given"
	}
	args := strings.Split(message.Message[12:], " ")
	channel := args[0]

	resp, err := GetStreamInfo(channel)
	if err != nil {
		log.Println(err)
	}
	if len(resp.Data.Streams) < 1 {
		return "Not live"
	}

	result := ""
	result += resp.Data.Streams[0].Title + ", "
	result += resp.Data.Streams[0].GameName + ", "
	result += strconv.Itoa(resp.Data.Streams[0].ViewerCount)
	return result
}

func (cmdHandler *CommandHandler) HttpStatusCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 12 {
		return "No code provided"
	}
	args := strings.Split(message.Message[12:], " ")

	code, err := strconv.Atoi(args[0])
	if err != nil {
		return "No valid number provided"
	}
	result := http.StatusText(code)
	return result
}

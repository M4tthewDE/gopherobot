package main

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type Bot struct {
	config    Config
	client    *twitch.Client
	startTime time.Time
	channels  []string
}

func NewBot() *Bot {
	bot := Bot{}
	bot.config = GetConfig()

	bot.client = twitch.NewClient("gopherobot", "oauth:"+bot.config.Twitch.Token)
	return &bot
}

func (b *Bot) Run() {
	b.startTime = time.Now()

	b.client.Join(b.config.Bot.Channels...)
	b.channels = append(b.channels, b.config.Bot.Channels...)

	b.client.OnPrivateMessage(b.onMessage)
	b.client.OnWhisperMessage(b.onWhisper)
}

func (b *Bot) onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if prefix == b.config.Bot.Prefix && message.User.ID == "116672490" {
		b.doCommand(message)
	}
}

func (b *Bot) doCommand(message twitch.PrivateMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		b.client.Say(message.Channel, EchoCommand(message))
	case "id":
		b.client.Say(message.Channel, UserIdCommand(message))
	case "user":
		b.client.Say(message.Channel, UserCommand(message))
	case "addfollowalert":
		b.client.Say(message.Channel, AddFollowAlertCommand(message))
	case "removefollowalert":
		b.client.Say(message.Channel, RemoveFollowAlertCommand(message))
	case "getfollowalerts":
		b.client.Say(message.Channel, GetFollowAlertsCommand())
	case "ping":
		b.client.Say(message.Channel, PingCommand(message))
	case "rawmsg":
		b.client.Say(message.Channel, RawMsgCommand(message.Raw))
	case "tmpjoin":
		b.client.Say(message.Channel, TmpJoinCommand(message))
	case "tmpleave":
		b.client.Say(message.Channel, TmpLeaveCommand(message))
	case "getchannels":
		b.client.Say(message.Channel, GetChannelsCommand(message))
	case "urlencode":
		b.client.Say(message.Channel, UrlEncodeCommand(message))
	case "urldecode":
		b.client.Say(message.Channel, UrlDecodeCommand(message))
	case "streaminfo":
		b.client.Say(message.Channel, StreamInfoCommand(message))
	case "httpstatus":
		b.client.Say(message.Channel, HttpStatusCommand(message))
	}
}

func (b *Bot) onWhisper(message twitch.WhisperMessage) {
	prefix := message.Message[0:1]

	if prefix == b.config.Bot.Prefix && message.User.ID == "116672490" {
		b.doWhisperCommand(message)
	}
}

func (b *Bot) doWhisperCommand(message twitch.WhisperMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "saveclip":
		SaveClipCommand(message)
	}
}

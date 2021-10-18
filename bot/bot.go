package bot

import (
	"strings"
	"time"

	"de.com.fdm/gopherobot/cmd"
	"de.com.fdm/gopherobot/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type Bot struct {
	config     *config.Config
	cmdHandler *cmd.CommandHandler
	client     *twitch.Client
	channels   []string
}

func NewBot(config *config.Config) *Bot {
	var bot Bot
	bot.config = config

	bot.client = twitch.NewClient("gopherobot", "oauth:"+config.Twitch.Token)

	// some Networks might block 6697
	bot.client.IrcAddress = "irc.chat.twitch.tv:443"
	bot.cmdHandler = cmd.NewCommandHandler(config, time.Now(), bot.client, &bot.channels)

	return &bot
}

func (b *Bot) Run() {
	b.client.OnPrivateMessage(b.onMessage)
	b.client.OnWhisperMessage(b.onWhisper)

	b.client.Join(b.config.Bot.Channels...)
	b.channels = append(b.channels, b.config.Bot.Channels...)

	if err := b.client.Connect(); err != nil {
		panic(err)
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
		b.cmdHandler.SaveClipCommand(message)
	case "revoke":
		b.client.Say("gopherobot", b.cmdHandler.RevokeAuthCommand(message))
	}
}

func (b *Bot) onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if message.User.ID == "82008718" && strings.HasPrefix(message.Message, "pajaS 🚨") {
		b.client.Say(message.Channel, "monkaS 🚨")
	}

	if prefix == b.config.Bot.Prefix && message.User.ID == "116672490" {
		b.doCommand(message)
	}
}

func (b *Bot) doCommand(message twitch.PrivateMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		b.client.Say(message.Channel, b.cmdHandler.EchoCommand(message))
	case "id":
		b.client.Say(message.Channel, b.cmdHandler.UserIDCommand(message))
	case "user":
		b.client.Say(message.Channel, b.cmdHandler.UserCommand(message))
	case "addfollowalert":
		b.client.Say(message.Channel, b.cmdHandler.AddFollowAlertCommand(message))
	case "removefollowalert":
		b.client.Say(message.Channel, b.cmdHandler.RemoveFollowAlertCommand(message))
	case "getfollowalerts":
		b.client.Say(message.Channel, b.cmdHandler.GetFollowAlertsCommand())
	case "ping":
		b.client.Say(message.Channel, b.cmdHandler.PingCommand(message))
	case "rawmsg":
		b.client.Say(message.Channel, b.cmdHandler.RawMsgCommand(message.Raw))
	case "tmpjoin":
		b.client.Say(message.Channel, b.cmdHandler.TmpJoinCommand(message))
	case "tmpleave":
		b.client.Say(message.Channel, b.cmdHandler.TmpLeaveCommand(message))
	case "getchannels":
		b.client.Say(message.Channel, b.cmdHandler.GetChannelsCommand(message))
	case "urlencode":
		b.client.Say(message.Channel, b.cmdHandler.URLEncodeCommand(message))
	case "urldecode":
		b.client.Say(message.Channel, b.cmdHandler.URLDecodeCommand(message))
	case "streaminfo":
		b.client.Say(message.Channel, b.cmdHandler.StreamInfoCommand(message))
	case "httpstatus":
		b.client.Say(message.Channel, b.cmdHandler.HTTPStatusCommand(message))
	case "nextlaunch":
		b.client.Say(message.Channel, b.cmdHandler.NextLaunchCommand(message))
	}
}

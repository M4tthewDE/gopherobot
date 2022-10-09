package bot

import (
	"strings"
	"time"

	"de.com.fdm/gopherobot/commands"
	"de.com.fdm/gopherobot/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type Bot struct {
	config        *config.Config
	client        *twitch.Client
	pingClient    *twitch.Client
	latencyReader *LatencyReader
	channels      []string
	lastPing      time.Time
	startTime     time.Time
}

func NewBot(config *config.Config) *Bot {
	var bot Bot
	bot.config = config

	bot.client = twitch.NewClient("gopherobot", "oauth:"+config.Twitch.Token)
	bot.pingClient = twitch.NewAnonymousClient()

	// some Networks might block 6697
	bot.client.IrcAddress = "irc.chat.twitch.tv:443"
	bot.pingClient.IrcAddress = "irc.chat.twitch.tv:443"

	// lower PING interval for latency checking
	bot.pingClient.IdlePingInterval = 1 * time.Second

	bot.startTime = time.Now()

	bot.latencyReader = NewLatencyReader()

	return &bot
}

func (b *Bot) Run() {
	go b.latencyReader.Read()
	go b.RunPingService()

	b.client.OnPrivateMessage(b.onMessage)
	b.pingClient.OnPongMessage(b.onPong)
	b.pingClient.OnPingSent(b.onPingSent)

	b.client.Join(b.config.Bot.Channels...)
	b.channels = append(b.channels, b.config.Bot.Channels...)

	if err := b.client.Connect(); err != nil {
		panic(err)
	}
}

func (b *Bot) RunPingService() {
	if err := b.pingClient.Connect(); err != nil {
		panic(err)
	}
}

func (b *Bot) onPong(message twitch.PongMessage) {
	b.latencyReader.LatencyChannel <- time.Since(b.lastPing)
}

func (b *Bot) onPingSent() {
	b.lastPing = time.Now()
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
		b.client.Say(message.Channel, commands.Echo(message))
	case "ping":
		b.client.Say(message.Channel, commands.Ping(b.startTime, b.latencyReader.latency, b.config))
	}
}

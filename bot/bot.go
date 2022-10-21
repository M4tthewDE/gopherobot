package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	b.client.OnConnect(func() {
		for _, channel := range b.channels {
			b.sendMessage(channel, "Restarted...")
		}
	})

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
	commandDeadline := time.Duration(b.config.Bot.Timeout) * time.Millisecond
	prefix := message.Message[0:1]

	if prefix != b.config.Bot.Prefix || message.User.Name == b.config.Bot.Name {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandDeadline)

	defer cancel()

	result, err := b.doCommand(ctx, message)
	if err != nil {

		if errors.Is(err, context.DeadlineExceeded) {
			b.sendMessage(message.Channel, "Command execution deadline exceeded")
			return
		}

		if errors.Is(err, ErrCommandNotAllowed) {
			return
		}

		log.Println(err)
		b.sendMessage(message.Channel, "Error during command execution")

		return
	}

	b.sendMessage(message.Channel, result)
}

func (b *Bot) sendMessage(channel string, message string) {
	if b.config.Bot.Profile == "PROD" {
		b.client.Say(channel, message)
	} else {
		b.client.Say(channel, fmt.Sprintf("[%s] %s", b.config.Bot.Profile, message))
	}
}

var ErrCommandNotAllowed = errors.New("command not allowed")

func (b *Bot) doCommand(ctx context.Context, message twitch.PrivateMessage) (string, error) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		if message.User.Name != b.config.Bot.Owner {
			return "", ErrCommandNotAllowed
		}

		return commands.Echo(message), nil
	case "ping":
		return commands.Ping(b.startTime, b.latencyReader.latency, b.config), nil
	case "improveemote":
		result, err := commands.ImproveEmote(ctx, message)
		if err != nil {
			return "", fmt.Errorf("improve emote error: %w", err)
		}

		return result, nil
	}

	return "Command not found", nil
}

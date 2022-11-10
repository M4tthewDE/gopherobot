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
	firstConnect  bool
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
	bot.firstConnect = true

	return &bot
}

func (b *Bot) Run() {
	go b.latencyReader.Read()
	go b.RunPingService()

	b.client.OnPrivateMessage(b.onMessage)
	b.pingClient.OnPongMessage(b.onPong)
	b.pingClient.OnPingSent(b.onPingSent)

	b.client.OnConnect(func() {
		if b.firstConnect {
			for _, channel := range b.channels {
				b.sendMessage(channel, "Restarted...")
			}

			b.firstConnect = false
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

func (b *Bot) onMessage(msg twitch.PrivateMessage) {
	commandDeadline := time.Duration(b.config.Bot.Timeout) * time.Millisecond
	prefix := msg.Message[0:1]

	if prefix != b.config.Bot.Prefix || msg.User.Name == b.config.Bot.Name {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandDeadline)

	defer cancel()

	result, err := b.doCommand(ctx, msg)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			response := b.buildResponse(msg, "Command execution deadline exceeded", true)

			b.sendMessage(msg.Channel, response)

			return
		}

		if errors.Is(err, ErrCommandNotAllowed) {
			return
		}

		log.Println(err)

		response := b.buildResponse(msg, "Error during command execution", true)
		b.sendMessage(msg.Channel, response)

		return
	}

	b.sendMessage(msg.Channel, result)
}

func (b *Bot) buildResponse(msg twitch.PrivateMessage, content string, doesPing bool) string {
	if doesPing {
		content = fmt.Sprintf("%s, %s", msg.User.Name, content)
	}

	if b.config.Bot.Profile != "PROD" {
		content = fmt.Sprintf("[%s] %s", b.config.Bot.Profile, content)
	}

	return content
}

func (b *Bot) sendMessage(channel string, message string) {
	b.client.Say(channel, message)
}

var ErrCommandNotAllowed = errors.New("command not allowed")

// Commands decide how their return values look like.
func (b *Bot) doCommand(ctx context.Context, msg twitch.PrivateMessage) (string, error) {
	identifier := strings.Split(msg.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		if msg.User.Name != b.config.Bot.Owner {
			return "", ErrCommandNotAllowed
		}

		response := commands.Echo(msg)

		return b.buildResponse(msg, response, false), nil
	case "ping":
		response := commands.Ping(b.startTime, b.latencyReader.latency, b.config)

		return b.buildResponse(msg, response, true), nil
	case "improveemote":
		response, err := commands.ImproveEmote(ctx, msg)
		if err != nil {
			return "", fmt.Errorf("improve emote error: %w", err)
		}

		return b.buildResponse(msg, response, true), nil
	}

	return "Command not found", nil
}

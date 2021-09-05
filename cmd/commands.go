package cmd

import (
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"de.com.fdm/gopherobot/config"
	"de.com.fdm/gopherobot/provider"
	"de.com.fdm/gopherobot/util"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/hako/durafmt"
)

type CommandHandler struct {
	config         *config.Config
	twitchProvider *provider.TwitchProvider
	hasteProvider  *provider.HasteProvider
	fdmProvider    *provider.FeelsdankmanProvider
	spaceXProvider *provider.SpaceXProvider
	startTime      time.Time
	client         *twitch.Client
	channels       *[]string
}

func NewCommandHandler(config *config.Config,
	startTime time.Time,
	client *twitch.Client,
	channels *[]string,
) *CommandHandler {
	cmdHandler := CommandHandler{
		config:         config,
		twitchProvider: &provider.TwitchProvider{Config: config},
		hasteProvider:  &provider.HasteProvider{Config: config},
		fdmProvider:    &provider.FeelsdankmanProvider{Config: config},
		spaceXProvider: &provider.SpaceXProvider{},
		startTime:      startTime,
		client:         client,
		channels:       channels,
	}

	return &cmdHandler
}

func (c *CommandHandler) EchoCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 6 {
		return ""
	}

	return message.Message[6:]
}

func (c *CommandHandler) UserIDCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 4 {
		return `No user provided`
	}

	args := strings.Split(message.Message[4:], " ")

	id, err := c.twitchProvider.GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	}

	return "User-ID of " + args[0] + " is " + id
}

func (c *CommandHandler) UserCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 6 {
		return `No ID provided`
	}

	args := strings.Split(message.Message[6:], " ")

	id := args[0]

	user, err := c.twitchProvider.GetUser(id)
	if err != nil {
		return `Couldn't find User for "` + args[0] + `"`
	}

	return "Username for " + args[0] + " is " + user
}

func (c *CommandHandler) AddFollowAlertCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[16:], " ")

	id, err := c.twitchProvider.GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	}

	err = c.fdmProvider.RegisterWebhook(id, message.Channel, message.User.Name)
	if err != nil {
		return "Error adding follow alert!"
	}

	return "Added follow alert for " + args[0] + "!"
}

func (c *CommandHandler) RemoveFollowAlertCommand(message twitch.PrivateMessage) string {
	args := strings.Split(message.Message[19:], " ")

	id, err := c.twitchProvider.GetUserID(args[0])
	if err != nil {
		return `Couldn't find User-ID for "` + args[0] + `"`
	}

	err = c.fdmProvider.RemoveWebhook(id, message.User.Name, message.Channel)
	if err != nil {
		return "Error removing follow alert!"
	}

	return "Removed follow alert for " + args[0] + "!"
}

func (c *CommandHandler) GetFollowAlertsCommand() string {
	followWebhook, err := c.fdmProvider.GetWebhooks()
	if err != nil {
		return "Error getting Followalerts"
	}

	payload := "Total alerts: " + strconv.Itoa(len(followWebhook.Data)) + " Channels: "

	users := make([]string, 0)

	for _, webhook := range followWebhook.Data {
		id := webhook.Condition.BroadcasterUserID
		user, _ := c.twitchProvider.GetUser(id)
		users = append(users, user)
	}

	for _, user := range users {
		payload = payload + user + " "
	}

	return payload
}

func (c *CommandHandler) PingCommand(message twitch.PrivateMessage) string {
	uptime := time.Since(c.startTime)
	result := "Pong! Uptime: " + durafmt.Parse(uptime).LimitFirstN(2).String() + ","

	apiUptime, err := c.fdmProvider.GetAPIUptime()
	if err != nil {
		result += " API-Uptime: Unavailable monkaS"
	} else {
		result += " API-Uptime: " + apiUptime + ","
	}

	result += " Commit: " + c.config.Git.Commit + ","
	result += " Branch: " + c.config.Git.Branch + ","

	latency := util.GetLatency(message)
	result += " Latency to tmi: " + latency

	return result
}

func (c *CommandHandler) RawMsgCommand(rawMessage string) string {
	return c.hasteProvider.UploadToHaste(rawMessage)
}

func (c *CommandHandler) TmpJoinCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 9 {
		return `No channel provided`
	}

	args := strings.Split(message.Message[9:], " ")
	channel := args[0]

	c.client.Join(channel)
	*c.channels = append(*c.channels, channel)

	return "Joined #" + channel
}

func (c *CommandHandler) TmpLeaveCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 10 {
		return `No channel provided`
	}

	args := strings.Split(message.Message[10:], " ")
	channel := args[0]

	c.client.Depart(channel)

	var index int

	for i, c := range *c.channels {
		if c == channel {
			index = i
		}
	}

	*c.channels = append((*c.channels)[:index], (*c.channels)[index+1:]...)

	return "Left #" + channel
}

func (c *CommandHandler) GetChannelsCommand(message twitch.PrivateMessage) string {
	if message.Channel != "matthewde" && message.Channel != "gopherobot" {
		return "Command not available in this channel to prevent pings"
	}

	var result string
	for _, channel := range *c.channels {
		result = result + ", " + channel
	}

	result = "Joined Channels: [" + result[2:] + "]"

	return result
}

func (c *CommandHandler) URLEncodeCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 11 {
		return "Nothing to encode"
	}

	content := message.Message[11:]

	return url.QueryEscape(content)
}

func (c *CommandHandler) URLDecodeCommand(message twitch.PrivateMessage) string {
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

func (c *CommandHandler) StreamInfoCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 12 {
		return "No channel given"
	}

	args := strings.Split(message.Message[12:], " ")
	channel := args[0]

	resp, err := c.twitchProvider.GetStreamInfo(channel)
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

func (c *CommandHandler) HTTPStatusCommand(message twitch.PrivateMessage) string {
	if len(message.Message) < 12 {
		return "No code provided"
	}

	args := strings.Split(message.Message[12:], " ")

	code, err := strconv.Atoi(args[0])
	if err != nil {
		return "No valid code provided"
	}

	result := http.StatusText(code)

	return result
}

func (c *CommandHandler) NextLaunchCommand(message twitch.PrivateMessage) string {
	nextLaunch, err := c.spaceXProvider.GetNextLaunch()
	if err != nil {
		return err.Error()
	}

	result := nextLaunch.DateUtc.String() + " (UTC) | "
	result += "Name: " + nextLaunch.Name + " | "

	if len(result+"Details: "+nextLaunch.Details) > 500 {
		remaining := 500 - len(result) - len("Details: ") - 3
		result += "Details: " + nextLaunch.Details[:remaining] + "..."

		return result
	}

	result += "Details: " + nextLaunch.Details

	return result
}

func (c *CommandHandler) RevokeAuthCommand(message twitch.WhisperMessage) string {
	if len(message.Message) < 8 {
		return "No auth provided"
	}

	args := strings.Split(message.Message[8:], " ")

	err := c.twitchProvider.RevokeAuth(args[0])
	if err != nil {
		return err.Error()
	}

	return "Success!"
}

func (c *CommandHandler) SaveClipCommand(message twitch.WhisperMessage) {
	clipURL := strings.Split(message.Message[4:], " ")[1]
	cmd := exec.Command("youtube-dl", "-o", "../data/%(title)s.%(ext)s", clipURL)

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}

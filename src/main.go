package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"gopkg.in/yaml.v2"
)

var client *twitch.Client
var StartTime time.Time
var Conf Config

func main() {
	f, err := os.Open("../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Conf)
	if err != nil {
		log.Fatal(err)
	}

	StartTime = time.Now()
	client = twitch.NewClient("gopherobot", "oauth:"+Conf.Twitch.Token)

	client.OnPrivateMessage(onMessage)
	client.OnWhisperMessage(onWhisper)

	client.Join(Conf.Bot.Channels...)

	go remoteMessageHandler(client)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

func onWhisper(message twitch.WhisperMessage) {
	prefix := message.Message[0:1]

	if prefix == Conf.Bot.Prefix && message.User.ID == "116672490" {
		doWhisperCommand(message)
	}
}

func doWhisperCommand(message twitch.WhisperMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "saveclip":
		SaveClipCommand(message)
	}
}

func onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if prefix == Conf.Bot.Prefix && message.User.ID == "116672490" {
		doCommand(message)
	}
}

func doCommand(message twitch.PrivateMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		client.Say(message.Channel, EchoCommand(message))
	case "id":
		client.Say(message.Channel, UserIdCommand(message))
	case "user":
		client.Say(message.Channel, UserCommand(message))
	case "addfollowalert":
		client.Say(message.Channel, AddFollowAlertCommand(message))
	case "removefollowalert":
		client.Say(message.Channel, RemoveFollowAlertCommand(message))
	case "getfollowalerts":
		client.Say(message.Channel, GetFollowAlertsCommand())
	case "ping":
		client.Say(message.Channel, PingCommand())
	case "rawmsg":
		client.Say(message.Channel, RawMsgCommand(message.Raw))
	}
}

func remoteMessageHandler(client *twitch.Client) {
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Println(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Println(err)
	}

	messageHandler := new(MessageHandler)
	err = rpc.Register(messageHandler)
	if err != nil {
		log.Println(err)
	}
	rpc.Accept(inbound)
}

type MessageHandler int

func (m *MessageHandler) SendMessage(content string, ack *bool) error {
	channel := strings.Split(content, " ")[0]
	message := content[len(channel):]

	client.Say(channel, message)
	return nil
}

type Config struct {
	Bot struct {
		Channels []string `yaml:"channels,flow"`
		Prefix   string   `yaml:"prefix"`
	} `yaml:"bot"`
	Api struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"api"`
	Twitch struct {
		Token     string `yaml:"token"`
		Client_ID string `yaml:"client_id"`
	} `yaml:"twitch"`
	Haste struct {
		Url string `yaml:"url"`
	} `yaml:"haste"`
}

package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"
)

var client *twitch.Client
var startTime time.Time
var config Config

func main() {
	f, err := os.Open("../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	startTime = time.Now()
	client = twitch.NewClient("gopherobot", "oauth:"+os.Getenv("TWITCH_TOKEN"))

	client.OnPrivateMessage(onMessage)

	//client.Join("gopherobot", "matthewde", "turtoise")
	client.Join(config.Bot.Channels...)

	go remoteMessageHandler(client)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

func onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if prefix == config.Bot.Prefix && message.User.ID == "116672490" {
		doCommand(message)
	}
}

func doCommand(message twitch.PrivateMessage) {
	identifier := strings.Split(message.Message, " ")[0][1:]
	switch identifier {
	case "echo":
		EchoCommand(message, client)
	case "id":
		UserIdCommand(message, client)
	case "addfollowalert":
		AddFollowAlertCommand(message, client, config.Api.Host)
	case "removefollowalert":
		RemoveFollowAlertCommand(message, client, config.Api.Host)
	case "ping":
		PingCommand(message, client, startTime)
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
	} `yaml:"api"`
}

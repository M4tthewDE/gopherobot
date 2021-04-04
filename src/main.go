package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
)

var client *twitch.Client

func main() {
	client = twitch.NewClient("gopherobot", "oauth:"+os.Getenv("TWITCH_TOKEN"))

	client.OnPrivateMessage(onMessage)

	client.Join("gopherobot", "matthewde", "turtoise")

	go remoteMessageHandler(client)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func onMessage(message twitch.PrivateMessage) {
	prefix := message.Message[0:1]

	if prefix == ";" && message.User.ID == "116672490" {
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
		AddFollowAlertCommand(message, client)
	case "removefollowalert":
		RemoveFollowAlertCommand(message, client)
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

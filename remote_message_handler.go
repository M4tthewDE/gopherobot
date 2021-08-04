package main

import (
	"log"
	"net"
	"net/rpc"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

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

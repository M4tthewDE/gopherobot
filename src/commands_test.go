package main

import (
	"github.com/gempir/go-twitch-irc/v2"
	"testing"
)

func TestEchoCommand(t *testing.T) {
	var msg twitch.PrivateMessage
	msg.Message = ";echo test"

	got := EchoCommand(msg)
	if got != "test" {
		t.Errorf("EchoCommand() = %s; want test", got)
	}
}

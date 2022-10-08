package util

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

func GetLatency(message twitch.PrivateMessage) string {
	latency := time.Since(message.Time)

	return latency.String()
}

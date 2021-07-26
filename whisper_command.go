package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

func SaveClipCommand(message twitch.WhisperMessage) {
	clip_url := strings.Split(message.Message[4:], " ")[1]

	cmd := exec.Command("youtube-dl", "-o", "../data/%(title)s.%(ext)s", clip_url)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

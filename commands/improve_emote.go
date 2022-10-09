package commands

import (
	"errors"
	"log"
	"strings"

	"de.com.fdm/gopherobot/providers"
	"github.com/gempir/go-twitch-irc/v2"
)

func ImproveEmote(message twitch.PrivateMessage) string {
	targetEmoteCode, err := getTargetEmoteCode(message)
	if err != nil {
		// TODO: this does not return anything
		return err.Error()
	}

	// get emotes for channel
	// TODO: add more emote providers
	emotes, err := providers.GetBttvEmotes(message.RoomID)
	if err != nil {
		return "Error getting emote"
	}

	// check channel bttv emotes
	for _, emote := range emotes.ChannelEmotes {
		if emote.Code == targetEmoteCode {
			result, err := improveBttvEmote(emote.Id)
			if err != nil {
				return err.Error()
			}

			return result
		}
	}

	// check shared bttv emotes
	for _, emote := range emotes.SharedEmotes {
		if emote.Code == targetEmoteCode {
			result, err := improveBttvEmote(emote.Id)
			if err != nil {
				return err.Error()
			}

			return result
		}
	}

	return "Improved"
}

func improveBttvEmote(emoteID string) (string, error) {
	log.Println(emoteID)
	// TODO: implement
	// do random transformations
	// https://golangdocs.com/golang-image-processing
	return "", nil
}

var errNoEmoteProvided = errors.New("no emote provided")

func getTargetEmoteCode(message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) < 2 {
		return "", errNoEmoteProvided
	}

	return words[1], nil
}

package commands

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"de.com.fdm/gopherobot/providers"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/h2non/bimg"
)

func ImproveEmote(message twitch.PrivateMessage) string {
	targetEmoteCode, err := getTargetEmoteCode(message.Message)
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
			_, err := improveBttvEmote(emote.Id)
			if err != nil {
				log.Println(err)
				return "Error improving emote"
			}

			return "DONE"
		}
	}

	// check shared bttv emotes
	for _, emote := range emotes.SharedEmotes {
		if emote.Code == targetEmoteCode {
			// FIXME: do something with result
			_, err := improveBttvEmote(emote.Id)
			if err != nil {
				log.Println(err)
				return "Error improving emote"
			}

			return "DONE"
		}
	}

	return "Improved"
}

func improveBttvEmote(emoteID string) ([]byte, error) {
	// TODO: implement and test with various emotes
	// do random transformations
	// https://golangdocs.com/golang-image-processing

	resp, err := http.Get("https://cdn.betterttv.net/emote/" + emoteID + "/3x")
	if err != nil {
		return nil, err
	}

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img := bimg.NewImage(buffer)

	newImage, err := img.Zoom(2)
	if err != nil {
		return nil, err
	}

	return newImage, nil

}

var errNoEmoteProvided = errors.New("no emote provided")

func getTargetEmoteCode(message string) (string, error) {
	message = strings.TrimSpace(message)
	words := strings.Split(message, " ")
	if len(words) < 2 {
		return "", errNoEmoteProvided
	}

	return words[1], nil
}

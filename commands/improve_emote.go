package commands

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"de.com.fdm/gopherobot/providers"
	vips "github.com/davidbyttow/govips/v2/vips"
	"github.com/gempir/go-twitch-irc/v2"
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
			newEmote, err := improveBttvEmote(emote.Id)
			if err != nil {
				log.Println(err)
				return "Error improving emote"
			}

			//TODO: upload instead of saving
			err = ioutil.WriteFile("new.gif", newEmote, 0644)
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

	vips.Startup(nil)
	defer vips.Shutdown()

	importParams := vips.NewImportParams()
	// needed to import all pages (frames)
	importParams.NumPages.Set(-1)

	image, err := vips.LoadImageFromBuffer(buffer, importParams)
	if err != nil {
		return nil, err
	}

	err = image.Flip(vips.DirectionHorizontal)
	if err != nil {
		return nil, err
	}

	pageDelays, err := image.PageDelay()
	if err != nil {
		return nil, err
	}

	// 2x the speed
	newPageDelays := make([]int, len(pageDelays))
	for i, delay := range pageDelays {
		newPageDelays[i] = delay / 4
	}

	image.SetPageDelay(newPageDelays)

	// widen emote
	err = image.ResizeWithVScale(2, 1, vips.KernelLanczos3)
	if err != nil {
		return nil, err
	}

	modifiedBuffer, _, err := image.ExportNative()
	if err != nil {
		return nil, err
	}

	return modifiedBuffer, nil
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

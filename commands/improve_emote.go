package commands

import (
	"errors"
	"log"
	"strings"

	"de.com.fdm/gopherobot/providers"
	vips "github.com/davidbyttow/govips/v2/vips"
	"github.com/gempir/go-twitch-irc/v2"
)

var errMessage = "Error improving emote"

func ImproveEmote(message twitch.PrivateMessage) string {
	targetEmoteCode, err := getTargetEmoteCode(message.Message)
	if err != nil {
		log.Println(err)

		return ""
	}

	// check bttv first
	emoteBuffer, didFind, err := findBttvEmote(targetEmoteCode, message.RoomID)
	if err != nil {
		log.Println(err)
	}

	if didFind {
		newEmoteBuffer, err := modifyEmote(emoteBuffer)
		if err != nil {
			log.Println(err)

			return errMessage
		}

		url, err := providers.UploadToKappaLol(newEmoteBuffer)
		if err != nil {
			log.Println(err)

			return errMessage
		}

		return url
	}

	// check 7tv
	emoteBuffer, didFind, err = findSevenTvEmote(targetEmoteCode, message.RoomID)
	if err != nil {
		log.Println(err)
	}

	if didFind {
		newEmoteBuffer, err := modifyEmote(emoteBuffer)
		if err != nil {
			log.Println(err)

			return errMessage
		}

		url, err := providers.UploadToKappaLol(newEmoteBuffer)
		if err != nil {
			log.Println(err)

			return errMessage
		}

		return url
	}

	return "Emote not found"
}

var errFindingSevenTvEmote = errors.New("error finding 7tv emote")

func findSevenTvEmote(targetEmoteName string, roomID string) ([]byte, bool, error) {
	emotes, err := providers.GetSevenTvEmotes(roomID)
	if err != nil {
		return nil, false, errFindingSevenTvEmote
	}

	for _, emote := range emotes {
		if emote.Name == targetEmoteName {
			emoteBuffer, err := providers.GetSevenTvEmote(emote.ID)
			if err != nil {
				return nil, false, errFindingSevenTvEmote
			}

			return emoteBuffer, true, nil
		}
	}

	return nil, false, nil
}

var errFindingBttvEmote = errors.New("error finding bttv emote")

func findBttvEmote(targetEmoteCode string, roomID string) ([]byte, bool, error) {
	// get bttv emotes for channel
	emotes, err := providers.GetBttvEmotes(roomID)
	if err != nil {
		return nil, false, errFindingBttvEmote
	}

	// check channel bttv emotes
	for _, emote := range emotes.ChannelEmotes {
		if emote.Code == targetEmoteCode {
			emoteBuffer, err := providers.GetBttvEmote(emote.ID)
			if err != nil {
				return nil, false, errFindingBttvEmote
			}

			return emoteBuffer, true, nil
		}
	}

	// check shared bttv emotes
	for _, emote := range emotes.SharedEmotes {
		if emote.Code == targetEmoteCode {
			emoteBuffer, err := providers.GetBttvEmote(emote.ID)
			if err != nil {
				return nil, true, errFindingBttvEmote
			}

			return emoteBuffer, true, nil
		}
	}

	return nil, false, nil
}

var errModifyingEmote = errors.New("failed to modify emote")

func modifyEmote(emoteBuffer []byte) ([]byte, error) {
	importParams := vips.NewImportParams()
	// needed to import all pages (frames)
	importParams.NumPages.Set(-1)

	image, err := vips.LoadImageFromBuffer(emoteBuffer, importParams)
	if err != nil {
		log.Println(err)

		return nil, errModifyingEmote
	}

	pageDelays, err := image.PageDelay()
	if err != nil {
		log.Println(err)

		return nil, errModifyingEmote
	}

	// 2x the speed
	newPageDelays := make([]int, len(pageDelays))

	for index, delay := range pageDelays {
		// if the delay is 10 or lower, it actually slows it down
		newDelay := delay / 4
		if newDelay < 11 {
			newDelay = 11
		}

		newPageDelays[index] = newDelay
	}

	err = image.SetPageDelay(newPageDelays)
	if err != nil {
		log.Println(err)

		return nil, errModifyingEmote
	}

	// widen emote
	err = image.ResizeWithVScale(2, 1, vips.KernelLanczos3)
	if err != nil {
		log.Println(err)

		return nil, errModifyingEmote
	}

	modifiedBuffer, _, err := image.ExportNative()
	if err != nil {
		log.Println(err)

		return nil, errModifyingEmote
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

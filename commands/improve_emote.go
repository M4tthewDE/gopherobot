package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"de.com.fdm/gopherobot/providers"
	vips "github.com/davidbyttow/govips/v2/vips"
	"github.com/gempir/go-twitch-irc/v2"
)

var ErrEmoteNotFound = errors.New("emote not found")

func ImproveEmote(ctx context.Context, message twitch.PrivateMessage) (string, error) {
	targetEmoteCode, err := getTargetEmoteCode(message.Message)
	if err != nil {
		return "", fmt.Errorf("get target emote code error: %w", err)
	}

	// check bttv first
	emoteBuffer, didFind, err := findBttvEmote(ctx, targetEmoteCode, message.RoomID)
	if err != nil {
		return "", fmt.Errorf("find bttv emote error: %w", err)
	}

	if didFind {
		newEmoteBuffer, err := modifyEmote(emoteBuffer)
		if err != nil {
			return "", fmt.Errorf("modify emote error: %w", err)
		}

		url, err := providers.UploadToKappaLol(ctx, newEmoteBuffer)
		if err != nil {
			return "", fmt.Errorf("upload to kappa.lol error: %w", err)
		}

		return url, nil
	}

	// check 7tv
	emoteBuffer, didFind, err = findSevenTvEmote(ctx, targetEmoteCode, message.RoomID)
	if err != nil {
		return "", fmt.Errorf("find 7tv emote error: %w", err)
	}

	if didFind {
		newEmoteBuffer, err := modifyEmote(emoteBuffer)
		if err != nil {
			return "", fmt.Errorf("modify emote error: %w", err)
		}

		url, err := providers.UploadToKappaLol(ctx, newEmoteBuffer)
		if err != nil {
			return "", fmt.Errorf("upload to kappa.lol error: %w", err)
		}

		return url, nil
	}

	return "", ErrEmoteNotFound
}

func findSevenTvEmote(ctx context.Context, targetEmoteName string, roomID string) ([]byte, bool, error) {
	emotes, err := providers.GetSevenTvEmotes(ctx, roomID)
	if err != nil {
		return nil, false, fmt.Errorf("getting 7tv emotes error: %w", err)
	}

	for _, emote := range emotes {
		if emote.Name == targetEmoteName {
			emoteBuffer, err := providers.GetSevenTvEmote(ctx, emote.ID)
			if err != nil {
				return nil, false, fmt.Errorf("getting 7tv emote error: %w", err)
			}

			return emoteBuffer, true, nil
		}
	}

	return nil, false, nil
}

var ErrFindingBttvEmote = errors.New("error finding bttv emote")

func findBttvEmote(ctx context.Context, targetEmoteCode string, roomID string) ([]byte, bool, error) {
	// get bttv emotes for channel
	emotes, err := providers.GetBttvEmotes(ctx, roomID)
	if err != nil {
		return nil, false, fmt.Errorf("getting bttv emotes error: %w", err)
	}

	// check channel bttv emotes
	for _, emote := range emotes.ChannelEmotes {
		if emote.Code == targetEmoteCode {
			emoteBuffer, err := providers.GetBttvEmote(ctx, emote.ID)
			if err != nil {
				return nil, false, fmt.Errorf("finding bttv emote error: %w", err)
			}

			return emoteBuffer, true, nil
		}
	}

	// check shared bttv emotes
	for _, emote := range emotes.SharedEmotes {
		if emote.Code == targetEmoteCode {
			emoteBuffer, err := providers.GetBttvEmote(ctx, emote.ID)
			if err != nil {
				return nil, false, fmt.Errorf("finding bttv emote error: %w", err)
			}

			return emoteBuffer, true, nil
		}
	}

	return nil, false, nil
}

func modifyEmote(emoteBuffer []byte) ([]byte, error) {
	importParams := vips.NewImportParams()
	// needed to import all pages (frames)
	importParams.NumPages.Set(-1)

	image, err := vips.LoadImageFromBuffer(emoteBuffer, importParams)
	if err != nil {
		return nil, fmt.Errorf("load image error: %w", err)
	}

	pageDelays, err := image.PageDelay()
	if err != nil {
		return nil, fmt.Errorf("get page delay error: %w", err)
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
		return nil, fmt.Errorf("set page delay error: %w", err)
	}

	// widen emote
	err = image.ResizeWithVScale(2, 1, vips.KernelLanczos3)
	if err != nil {
		return nil, fmt.Errorf("resize error: %w", err)
	}

	modifiedBuffer, _, err := image.ExportNative()
	if err != nil {
		return nil, fmt.Errorf("export error: %w", err)
	}

	return modifiedBuffer, nil
}

var ErrNoEmoteProvided = errors.New("no emote provided")

func getTargetEmoteCode(message string) (string, error) {
	message = strings.TrimSpace(message)

	words := strings.Split(message, " ")
	if len(words) < 2 {
		return "", ErrNoEmoteProvided
	}

	return words[1], nil
}

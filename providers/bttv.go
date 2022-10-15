package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetBttvEmotes(ctx context.Context, userID string) (*BttvEmotes, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.betterttv.net/3/cached/users/twitch/"+userID,
		nil)
	if err != nil {
		return nil, fmt.Errorf("request build error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	defer resp.Body.Close()

	var emotes BttvEmotes

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	return &emotes, nil
}

func GetBttvEmote(ctx context.Context, emoteID string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://cdn.betterttv.net/emote/"+emoteID+"/3x",
		nil)
	if err != nil {
		return nil, fmt.Errorf("request build error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	defer resp.Body.Close()

	emoteBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("emote read error: %w", err)
	}

	return emoteBuffer, nil
}

type BttvEmotes struct {
	ChannelEmotes []ChannelEmote `json:"channelEmotes"`
	SharedEmotes  []SharedEmote  `json:"sharedEmotes"`
}

type ChannelEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type SharedEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

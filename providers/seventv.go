package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetSevenTvEmotes(ctx context.Context, userID string) ([]SevenTvEmote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.7tv.app/v2/users/%s/emotes", userID),
		nil)
	if err != nil {
		return nil, fmt.Errorf("request build error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	defer resp.Body.Close()

	var emotes []SevenTvEmote

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	return emotes, nil
}

func GetSevenTvEmote(ctx context.Context, id string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x", id),
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
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	return emoteBuffer, nil
}

type SevenTvEmote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

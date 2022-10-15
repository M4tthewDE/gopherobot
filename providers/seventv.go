package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var ErrFetchingSevenTvEmotes = errors.New("error fetching 7tv emotes")

func GetSevenTvEmotes(userID string) ([]SevenTvEmote, error) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("https://api.7tv.app/v2/users/%s/emotes", userID),
		nil)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	defer resp.Body.Close()

	var emotes []SevenTvEmote

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	return emotes, nil
}

func GetSevenTvEmote(id string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x", id),
		nil)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	defer resp.Body.Close()

	emoteBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil, ErrFetchingSevenTvEmotes
	}

	return emoteBuffer, nil
}

type SevenTvEmote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

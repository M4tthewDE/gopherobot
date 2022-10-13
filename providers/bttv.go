package providers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var errFetchingBttvEmotes = errors.New("error fetching bttv emotes")

func GetBttvEmotes(userID string) (*BttvEmotes, error) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		"https://api.betterttv.net/3/cached/users/twitch/"+userID,
		nil)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
	}

	req.Header.Add("Content-Type", "multipart/form-data")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
	}

	defer resp.Body.Close()

	var emotes BttvEmotes

	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
	}

	return &emotes, nil
}

func GetBttvEmote(emoteID string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		"https://cdn.betterttv.net/emote/"+emoteID+"/3x",
		nil)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
	}

	defer resp.Body.Close()

	emoteBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return nil, errFetchingBttvEmotes
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

package providers

import (
	"encoding/json"
	"net/http"
)

func GetBttvEmotes(userID string) (*BttvEmotes, error) {
	resp, err := http.Get("https://api.betterttv.net/3/cached/users/twitch/" + userID)
	if err != nil {
		return nil, err
	}

	var emotes BttvEmotes
	err = json.NewDecoder(resp.Body).Decode(&emotes)
	if err != nil {
		return nil, err
	}

	return &emotes, nil
}

type BttvEmotes struct {
	ChannelEmotes []ChannelEmote `json:"channelEmotes"`
	SharedEmotes  []SharedEmote  `json:"sharedEmotes"`
}

type ChannelEmote struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type SharedEmote struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

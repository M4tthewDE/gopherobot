package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nicklaw5/helix"
)

func GetUserID(user string) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        Conf.Twitch.Client_ID,
		UserAccessToken: Conf.Twitch.Token,
	})
	if err != nil {
		return "", err
	}
	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{user},
	})
	if err != nil {
		return "", err
	}
	return resp.Data.Users[0].ID, nil
}

func GetUser(id int) (string, error) {
	url := "https://api.twitch.tv/helix/users?id=" + strconv.Itoa(id)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Client-ID", Conf.Twitch.Client_ID)
	req.Header.Add("Authorization", "Bearer "+Conf.Twitch.Token)

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	if r.StatusCode != 200 {
		return "", errors.New(strconv.Itoa(r.StatusCode))
	}

	var getUserIdJson GetUserIdJson
	err = json.Unmarshal(data, &getUserIdJson)
	if err != nil {
		return "", err
	}

	if len(getUserIdJson.Data) == 0 {
		return "", errors.New("no User-ID found")
	}

	return getUserIdJson.Data[0].DisplayName, nil
}

type GetUserIdJson struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

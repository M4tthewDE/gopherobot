package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetUserID(user string) (int, error) {
	url := "https://api.twitch.tv/helix/users?login=" + user
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Client-ID", Conf.Twitch.Client_ID)
	req.Header.Add("Authorization", "Bearer "+Conf.Twitch.Token)

	r, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	if r.StatusCode != 200 {
		return 0, errors.New(strconv.Itoa(r.StatusCode))
	}

	var getUserIdJson GetUserIdJson
	err = json.Unmarshal(data, &getUserIdJson)
	if err != nil {
		return 0, err
	}

	if len(getUserIdJson.Data) == 0 {
		return 0, errors.New("no User-ID found")
	}

	id, err := strconv.Atoi(getUserIdJson.Data[0].ID)
	if err != nil {
		return 0, err
	}
	return id, nil
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

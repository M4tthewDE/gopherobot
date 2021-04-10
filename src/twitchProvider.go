package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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

	req.Header.Add("Client-ID", os.Getenv("TWITCH_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TWITCH_TOKEN"))

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
		return 0, errors.New("No User-ID found!")
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

	req.Header.Add("Client-ID", os.Getenv("TWITCH_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TWITCH_TOKEN"))

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
		return "", errors.New("No User-ID found!")
	}

	return getUserIdJson.Data[0].DisplayName, nil
}

func GetActiveSubscriptions() (ActiveSubscriptions, error) {
	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Client-ID", os.Getenv("TWITCH_APP_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TWITCH_APP_TOKEN"))

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return ActiveSubscriptions{}, err
	}
	var activeSubs ActiveSubscriptions
	err = json.Unmarshal(data, &activeSubs)
	if err != nil {
		return ActiveSubscriptions{}, err
	}

	if r.StatusCode != 200 {
		return ActiveSubscriptions{}, errors.New(strconv.Itoa(r.StatusCode))
	}

	return activeSubs, nil
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

type ActiveSubscriptions struct {
	Total int `json:"total"`
	Data  []struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		CreatedAt time.Time `json:"created_at"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		Cost int `json:"cost"`
	} `json:"data"`
	Limit        int `json:"limit"`
	MaxTotalCost int `json:"max_total_cost"`
	TotalCost    int `json:"total_cost"`
	Pagination   struct {
	} `json:"pagination"`
}

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

func RegisterWebhook(id int, channel string, name string) error {
	url := "https://" + Conf.Api.Host + "/webhook/register?type=follow&id=" + strconv.Itoa(id) + "&user=" + name + "&channel=" + channel
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New(strconv.Itoa(r.StatusCode))
	}
	return nil
}

func RemoveWebhook(id int) error {
	activeSubs, err := GetActiveSubscriptions()
	if err != nil {
		log.Println(err)
	}
	for _, sub := range activeSubs.Data {
		if sub.Condition.BroadcasterUserID == strconv.Itoa(id) {
			err = DeleteWebhook(sub.ID)
			if err != nil {
				log.Println(err)
			}
			return nil
		}
	}
	return errors.New("No webhook for this user found!")
}

func DeleteWebhook(id string) error {
	url := "https://" + Conf.Api.Host + "/webhook/twitch/setup/delete?id=" + id
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New(strconv.Itoa(r.StatusCode))
	}
	return nil
}

func GetWebhooks() (FollowWebhook, error) {
	url := "https://" + Conf.Api.Host + "/webhook/twitch/setup/subscriptions"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FollowWebhook{}, err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return FollowWebhook{}, err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return FollowWebhook{}, err
	}

	var followWebhook FollowWebhook
	err = json.Unmarshal(data, &followWebhook)
	if err != nil {
		return FollowWebhook{}, err
	}

	if r.StatusCode != 200 {
		return FollowWebhook{}, errors.New(strconv.Itoa(r.StatusCode))
	}
	return followWebhook, nil
}

func GetApiUptime() (string, error) {
	url := "https://" + Conf.Api.Host + "/webhook/twitch/setup/uptime"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

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
	return string(data), nil
}

type FollowWebhook struct {
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

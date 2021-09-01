package provider

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"de.com.fdm/gopherobot/config"
)

type FeelsdankmanProvider struct {
	Config *config.Config
}

func (f *FeelsdankmanProvider) RegisterWebhook(id string, channel string, name string) error {
	url := "https://" + f.Config.API.Host + "/webhook/register?type=follow&id=" + id + "&user=" + name + "&channel=" + channel
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
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

func (f *FeelsdankmanProvider) RemoveWebhook(id string, username string, channel string) error {
	broadcaster_id := id
	url := "https://" + f.Config.API.Host + "/webhook/twitch/setup/delete?broadcaster=" + broadcaster_id + "&user=" + username + "&channel=" + channel
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
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

func (f *FeelsdankmanProvider) GetWebhooks() (FollowWebhook, error) {
	url := "https://" + f.Config.API.Host + "/webhook/twitch/setup/subscriptions"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FollowWebhook{}, err
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
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

func (f *FeelsdankmanProvider) GetApiUptime() (string, error) {
	url := "https://" + f.Config.API.Host + "/webhook/twitch/setup/uptime"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
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

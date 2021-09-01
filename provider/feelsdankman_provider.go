package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"de.com.fdm/gopherobot/config"
)

var errBadStatus = errors.New("request returned bad statuscode")

type FeelsdankmanProvider struct {
	Config *config.Config
}

func (f *FeelsdankmanProvider) RegisterWebhook(id string, channel string, name string) error {
	url := "https://" + f.Config.API.Host
	url += "/webhook/register?type=follow&id=" + id + "&user=" + name + "&channel=" + channel
	client := &http.Client{}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error registering webhook: %w", err)
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error registering webhook: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errBadStatus
	}

	return nil
}

func (f *FeelsdankmanProvider) RemoveWebhook(id string, username string, channel string) error {
	broadcasterID := id
	url := "https://" + f.Config.API.Host
	url += "/webhook/twitch/setup/delete?broadcaster=" + broadcasterID + "&user=" + username + "&channel=" + channel
	client := &http.Client{}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error removing webhook: %w", err)
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error removing webhook: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errBadStatus
	}

	return nil
}

func (f *FeelsdankmanProvider) GetWebhooks() (FollowWebhook, error) {
	url := "https://" + f.Config.API.Host + "/webhook/twitch/setup/subscriptions"
	client := &http.Client{}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return FollowWebhook{}, fmt.Errorf("error getting webhooks: %w", err)
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return FollowWebhook{}, fmt.Errorf("error getting webhooks: %w", err)
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return FollowWebhook{}, fmt.Errorf("error getting webhooks: %w", err)
	}

	var followWebhook FollowWebhook

	err = json.Unmarshal(data, &followWebhook)
	if err != nil {
		return FollowWebhook{}, fmt.Errorf("error getting webhook: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return FollowWebhook{}, errBadStatus
	}

	return followWebhook, nil
}

func (f *FeelsdankmanProvider) GetAPIUptime() (string, error) {
	url := "https://" + f.Config.API.Host + "/webhook/twitch/setup/uptime"
	client := &http.Client{}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error uptime: %w", err)
	}

	user := f.Config.API.User
	pass := f.Config.API.Pass
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting uptime: %w", err)
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("error getting uptime: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return "", errBadStatus
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
	Limit        int      `json:"limit"`
	MaxTotalCost int      `json:"max_total_cost"`
	TotalCost    int      `json:"total_cost"`
	Pagination   struct{} `json:"pagination"`
}

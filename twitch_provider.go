package main

import (
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

func GetUser(id string) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        Conf.Twitch.Client_ID,
		UserAccessToken: Conf.Twitch.Token,
	})
	if err != nil {
		return "", err
	}
	resp, err := client.GetUsers(&helix.UsersParams{
		IDs: []string{id},
	})
	if err != nil {
		return "", err
	}
	return resp.Data.Users[0].Login, nil
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

package provider

import (
	"errors"
	"fmt"

	"de.com.fdm/gopherobot/config"
	"github.com/nicklaw5/helix"
)

var errUserNotFound = errors.New("no user found")

type TwitchProvider interface {
	GetUserID(user string) (string, error)
	GetUser(id string) (string, error)
	GetStreamInfo(user string) (*helix.StreamsResponse, error)
	RevokeAuth(auth string) error
}

type ActualTwitchProvider struct {
	Config *config.Config
}

func (t *ActualTwitchProvider) GetUserID(user string) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        t.Config.Twitch.ClientID,
		UserAccessToken: t.Config.Twitch.Token,
	})
	if err != nil {
		return "", fmt.Errorf("error getting user id: %w", err)
	}

	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{user},
	})
	if err != nil {
		return "", fmt.Errorf("error getting user id: %w", err)
	}

	if len(resp.Data.Users) == 0 {
		return "", errUserNotFound
	}

	return resp.Data.Users[0].ID, nil
}

func (t *ActualTwitchProvider) GetUser(id string) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        t.Config.Twitch.ClientID,
		UserAccessToken: t.Config.Twitch.Token,
	})
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	resp, err := client.GetUsers(&helix.UsersParams{
		IDs: []string{id},
	})
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	if len(resp.Data.Users) == 0 {
		return "", errUserNotFound
	}

	return resp.Data.Users[0].Login, nil
}

func (t *ActualTwitchProvider) GetStreamInfo(user string) (*helix.StreamsResponse, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        t.Config.Twitch.ClientID,
		UserAccessToken: t.Config.Twitch.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting stream: %w", err)
	}

	resp, err := client.GetStreams(&helix.StreamsParams{
		UserLogins: []string{user},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting stream: %w", err)
	}

	return resp, nil
}

func (t *ActualTwitchProvider) RevokeAuth(auth string) error {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        t.Config.Twitch.ClientID,
		UserAccessToken: t.Config.Twitch.Token,
	})
	if err != nil {
		return fmt.Errorf("error revoking auth: %w", err)
	}

	_, err = client.RevokeUserAccessToken(auth)
	if err != nil {
		return fmt.Errorf("error revoking auth: %w", err)
	}

	return nil
}

type TestTwitchProvider struct {
	Config *config.Config
}

func (t *TestTwitchProvider) GetUserID(user string) (string, error) {
	return "1337", nil
}

func (t *TestTwitchProvider) GetUser(id string) (string, error) {
	return "user", nil
}

func (t *TestTwitchProvider) GetStreamInfo(user string) (*helix.StreamsResponse, error) {
	response := helix.StreamsResponse{
		Data: helix.ManyStreams{
			Streams: []helix.Stream{
				{
					Title:       "test-title",
					GameName:    "test-game",
					ViewerCount: 1000,
				},
			},
		},
	}

	return &response, nil
}

func (t *TestTwitchProvider) RevokeAuth(auth string) error {
	return nil
}

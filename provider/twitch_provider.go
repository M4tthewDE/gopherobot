package provider

import (
	"errors"
	"fmt"

	"de.com.fdm/gopherobot/config"
	"github.com/nicklaw5/helix"
)

var errUserNotFound = errors.New("no user found")

type TwitchProvider struct {
	Config *config.Config
}

func (t *TwitchProvider) GetUserID(user string) (string, error) {
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

func (t *TwitchProvider) GetUser(userID string) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        t.Config.Twitch.ClientID,
		UserAccessToken: t.Config.Twitch.Token,
	})
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	resp, err := client.GetUsers(&helix.UsersParams{
		IDs: []string{userID},
	})
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	if len(resp.Data.Users) == 0 {
		return "", errUserNotFound
	}

	return resp.Data.Users[0].Login, nil
}

func (t *TwitchProvider) GetStreamInfo(user string) (*helix.StreamsResponse, error) {
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

func (t *TwitchProvider) RevokeAuth(auth string) error {
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

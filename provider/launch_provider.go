package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LaunchProvider interface {
	GetNextLaunch() (NextLaunch, error)
}

type SpaceXProvider struct{}

func (s SpaceXProvider) GetNextLaunch() (NextLaunch, error) {
	httpClient := &http.Client{}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.spacexdata.com/v5/launches/next", nil)
	if err != nil {
		return NextLaunch{}, fmt.Errorf("error getting next launch: %w", err)
	}

	// send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return NextLaunch{}, fmt.Errorf("error getting next launch: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return NextLaunch{}, fmt.Errorf("error getting next launch: %w", err)
	}

	var nextLaunch NextLaunch

	err = json.Unmarshal(data, &nextLaunch)
	if err != nil {
		return NextLaunch{}, fmt.Errorf("error getting next launch: %w", err)
	}

	return nextLaunch, nil
}

type TestLaunchProvider struct{}

func (t TestLaunchProvider) GetNextLaunch() (NextLaunch, error) {
	return NextLaunch{
		DateUtc: time.Date(2021, 1, 1, 12, 0, 0, 0, time.FixedZone("UTC", 0)),
		Name:    "Test-Launch",
		Details: "Test-Details",
	}, nil
}

type NextLaunch struct {
	Fairings interface{} `json:"fairings"`
	Links    struct {
		Patch struct {
			Small string `json:"small"`
			Large string `json:"large"`
		} `json:"patch"`
		Reddit struct {
			Campaign string      `json:"campaign"`
			Launch   string      `json:"launch"`
			Media    string      `json:"media"`
			Recovery interface{} `json:"recovery"`
		} `json:"reddit"`
		Flickr struct {
			Small    []interface{} `json:"small"`
			Original []string      `json:"original"`
		} `json:"flickr"`
		Presskit  string `json:"presskit"`
		Webcast   string `json:"webcast"`
		YoutubeID string `json:"youtube_id"`
		Article   string `json:"article"`
		Wikipedia string `json:"wikipedia"`
	} `json:"links"`
	StaticFireDateUtc  time.Time     `json:"static_fire_date_utc"`
	StaticFireDateUnix int           `json:"static_fire_date_unix"`
	Tdb                bool          `json:"tdb"`
	Net                bool          `json:"net"`
	Window             int           `json:"window"`
	Rocket             string        `json:"rocket"`
	Success            bool          `json:"success"`
	Failures           []interface{} `json:"failures"`
	Details            string        `json:"details"`
	Crew               []interface{} `json:"crew"`
	Ships              []interface{} `json:"ships"`
	Capsules           []string      `json:"capsules"`
	Payloads           []string      `json:"payloads"`
	Launchpad          string        `json:"launchpad"`
	AutoUpdate         bool          `json:"auto_update"`
	FlightNumber       int           `json:"flight_number"`
	Name               string        `json:"name"`
	DateUtc            time.Time     `json:"date_utc"`
	DateUnix           int           `json:"date_unix"`
	DateLocal          string        `json:"date_local"`
	DatePrecision      string        `json:"date_precision"`
	Upcoming           bool          `json:"upcoming"`
	Cores              []struct {
		Core           string `json:"core"`
		Flight         int    `json:"flight"`
		Gridfins       bool   `json:"gridfins"`
		Legs           bool   `json:"legs"`
		Reused         bool   `json:"reused"`
		LandingAttempt bool   `json:"landing_attempt"`
		LandingSuccess bool   `json:"landing_success"`
		LandingType    string `json:"landing_type"`
		Landpad        string `json:"landpad"`
	} `json:"cores"`
	ID string `json:"id"`
}

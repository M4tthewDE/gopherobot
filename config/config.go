package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	GitCommit string
	GitBranch string
)

func GetConfig() *Config {
	file, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(file)

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	config.Git.Commit = GitCommit
	config.Git.Branch = GitBranch

	return &config
}

type Config struct {
	Bot struct {
		Channels []string `yaml:"channels,flow"`
		Prefix   string   `yaml:"prefix"`
		Profile  string   `yaml:"profile"`
		Timeout  int      `yaml:"timeout"`
		Owner    string   `yaml:"owner"`
		Name     string   `yaml:"name"`
	} `yaml:"bot"`
	Twitch struct {
		Token    string `yaml:"token"`
		ClientID string `yaml:"clientId"`
	} `yaml:"twitch"`
	Git struct {
		Branch string
		Commit string
	}
}

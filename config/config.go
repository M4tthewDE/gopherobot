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
	} `yaml:"bot"`
	API struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"api"`
	Twitch struct {
		Token    string `yaml:"token"`
		ClientID string `yaml:"clientId"`
	} `yaml:"twitch"`
	Haste struct {
		URL string `yaml:"url"`
	} `yaml:"haste"`
	Git struct {
		Branch string
		Commit string
	}
}

package config

import (
	"log"
	"os"

	"de.com.fdm/gopherobot/util"
	"gopkg.in/yaml.v2"
)

func GetConfig() *Config {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	branch, err := util.GetBranch()
	if err != nil {
		log.Fatal(err)
	}

	config.Git.Branch = branch

	commit, err := util.GetCommit()
	if err != nil {
		log.Fatal(err)
	}

	config.Git.Commit = commit

	return &config
}

type Config struct {
	Bot struct {
		Channels []string `yaml:"channels,flow"`
		Prefix   string   `yaml:"prefix"`
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

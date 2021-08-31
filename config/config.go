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
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

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
	Api struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"api"`
	Twitch struct {
		Token     string `yaml:"token"`
		Client_ID string `yaml:"client_id"`
	} `yaml:"twitch"`
	Haste struct {
		Url string `yaml:"url"`
	} `yaml:"haste"`
	Git struct {
		Branch string
		Commit string
	}
}

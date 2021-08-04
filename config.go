package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

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

func GetConfig() Config {
	var config Config

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	branch, err := GetBranch()
	if err != nil {
		log.Fatal(err)
	}
	config.Git.Branch = branch

	commit, err := GetCommit()
	if err != nil {
		log.Fatal(err)
	}
	config.Git.Commit = commit
	return config
}

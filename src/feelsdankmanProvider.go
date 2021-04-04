package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

func RegisterWebhook(id int) error {
	url := "https://feelsdankman.xyz/webhook/register?type=follow&id=" + strconv.Itoa(id)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New(strconv.Itoa(r.StatusCode))
	}
	return nil
}

func RemoveWebhook(id int) error {
	activeSubs, err := GetActiveSubscriptions()
	if err != nil {
		log.Println(err)
	}
	for _, sub := range activeSubs.Data {
		if sub.Condition.BroadcasterUserID == strconv.Itoa(id) {
			err = DeleteWebhook(sub.ID)
			if err != nil {
				log.Println(err)
			}
			return nil
		}
	}
	return errors.New("No webhook for this user found!")
}

func DeleteWebhook(id string) error {
	url := "https://feelsdankman.xyz/webhook/twitch/setup/delete?id=" + id
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASS")
	req.SetBasicAuth(user, pass)

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New(strconv.Itoa(r.StatusCode))
	}
	return nil
}

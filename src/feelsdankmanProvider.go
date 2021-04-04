package main

import (
	"errors"
	"io"
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

func DeleteWebhook(id int) {
	log.Println(id)
}

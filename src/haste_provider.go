package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// inspired by https://github.com/zneix/haste-client/blob/master/main.go 🐍
func UploadToHaste(data string) string {

	type HasteResponseData struct {
		Key string `json:"key,omitempty"`
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest("POST", Conf.Haste.Url+"/documents", bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Println("New Request error: " + err.Error())
		return ""
	}

	//send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Request Do error: " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Println(fmt.Sprintf("Error while uploading data: %d", resp.StatusCode))
		return ""
	}

	//error out if the invite isn't found or something else went wrong with the request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Error while reading response: %s", err.Error()))
		return ""
	}

	var jsonResponse HasteResponseData
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Println(fmt.Sprintf("Error while unmarshaling JSON response: %s", err.Error()))
		return ""
	}

	var finalURL = Conf.Haste.Url
	finalURL += "/raw/" + jsonResponse.Key

	return finalURL
}

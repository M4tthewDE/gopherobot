package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"de.com.fdm/gopherobot/config"
)

type HasteProvider struct {
	Config *config.Config
}

// inspired by https://github.com/zneix/haste-client/blob/master/main.go
func (h *HasteProvider) UploadToHaste(data string) string {
	type HasteResponseData struct {
		Key string `json:"key,omitempty"`
	}

	httpClient := &http.Client{}
	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.Config.Haste.URL+"/documents",
		bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Println("New Request error: " + err.Error())

		return ""
	}

	// send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Request Do error: " + err.Error())

		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Printf("Error while uploading data: %d", resp.StatusCode)

		return ""
	}

	// error out if the invite isn't found or something else went wrong with the request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading response: %s", err.Error())

		return ""
	}

	var jsonResponse HasteResponseData
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Printf("Error while unmarshaling JSON response: %s", err.Error())

		return ""
	}

	finalURL := h.Config.Haste.URL
	finalURL += "/raw/" + jsonResponse.Key

	return finalURL
}

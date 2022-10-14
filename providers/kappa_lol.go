package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
)

var errUploadKappaLol = errors.New("failed to upload to kappa.lol")

func UploadToKappaLol(data []byte) (string, error) {
	body := bytes.NewBuffer(data)
	mp := multipart.NewWriter(body)

	_, err := mp.CreateFormFile("file", "doesThisEvenMatter.gif")
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	err = mp.Close()
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"https://kappa.lol/api/upload",
		body)
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	req.Header.Add("Content-Type", mp.FormDataContentType())

	log.Printf("%+v\n", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	defer resp.Body.Close()

	var kappaLolResponse KappaLolResponse

	err = json.NewDecoder(resp.Body).Decode(&kappaLolResponse)
	if err != nil {
		return "", errUploadKappaLol
	}

	return kappaLolResponse.Link, nil
}

type KappaLolResponse struct {
	Link string `json:"link"`
}

package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

var errUploadKappaLol = errors.New("failed to upload to kappa.lol")

func UploadToKappaLol(data []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "whoasked")
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	_, err = io.Copy(part, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	err = writer.Close()
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

	req.Header.Add("Content-Type", writer.FormDataContentType())

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

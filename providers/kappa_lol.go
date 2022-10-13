package providers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

var errUploadKappaLol = errors.New("failed to upload to kappa.lol")

func UploadToKappaLol(data []byte) (string, error) {
	body := bytes.NewBuffer(data)
	mp := multipart.NewWriter(body)

	part, err := mp.CreateFormFile("@file", "NewEmoteTEST123")
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	io.Copy(part, body)

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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	defer resp.Body.Close()

	urlBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

	log.Println(string(urlBuffer))

	return "TEST", nil

}

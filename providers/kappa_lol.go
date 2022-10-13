package providers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
)

var errUploadKappaLol = errors.New("failed to upload to kappa.lol")

func UploadToKappaLol(data []byte) (string, error) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"https://kappa.lol/api/upload",
		bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)

		return "", errUploadKappaLol
	}

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

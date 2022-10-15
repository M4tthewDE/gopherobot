package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadToKappaLol(ctx context.Context, data []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "whoasked")
	if err != nil {
		return "", fmt.Errorf("creating form error: %w", err)
	}

	_, err = io.Copy(part, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("form file copy error: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("writer close error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://kappa.lol/api/upload",
		body)
	if err != nil {
		return "", fmt.Errorf("request build error: %w", err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch error: %w", err)
	}

	defer resp.Body.Close()

	var kappaLolResponse KappaLolResponse

	err = json.NewDecoder(resp.Body).Decode(&kappaLolResponse)
	if err != nil {
		return "", fmt.Errorf("json decode error: %w", err)
	}

	return kappaLolResponse.Link, nil
}

type KappaLolResponse struct {
	Link string `json:"link"`
}

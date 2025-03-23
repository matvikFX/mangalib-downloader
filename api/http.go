package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
)

func ReqImg(ctx context.Context, pageURL string) ([]byte, error) {
	log := slog.With("API", "ReqImg")

	resp, err := req(ctx, pageURL)
	if err != nil {
		log.Error("Error getting response", "Error", err)
		return nil, err
	}
	defer resp.Body.Close()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading image", "Error", err)
		return nil, err
	}

	return img, nil
}

func ReqJSON(ctx context.Context, url string, data any) error {
	log := slog.With("API", "ReqJSON")

	resp, err := req(ctx, url)
	if err != nil {
		log.Error("Error getting response", "Error", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		log.Error("Error decondig response", "Error", err)
		return err
	}

	return nil
}

func req(ctx context.Context, url string) (*http.Response, error) {
	log := slog.With("API", "req")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error("Error creating request with context", "Error", err)
		return nil, err
	}

	var contType string
	switch filepath.Ext(url) {
	case ".jpg", ".jpeg", ".jpe", ".jif", ".jfif":
		contType = "image/jpeg"
	case ".gif":
		contType = "image/gif"
	case ".png":
		contType = "image/png"
	default:
		contType = "application/json"
	}

	header := http.Header{}
	header.Set("Content-Type", contType)
	header.Set("Authorization", readAuthToken())
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Error while doing a request", "Error", err)
		return nil, err
	}
	// log.Debug("Response status", "status", resp.Status)

	return resp, nil
}

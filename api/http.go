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
	url := createPageURL(pageURL)
	resp, err := req(ctx, url)
	if err != nil {
		slog.Error("Error getting response", "Error", err)
		return nil, err
	}
	defer resp.Body.Close()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading image", "Error", err)
		return nil, err
	}

	return img, nil
}

func ReqJSON(ctx context.Context, url string, data any) error {
	resp, err := req(ctx, url)
	if err != nil {
		slog.Error("Error getting response", "Error", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		slog.Error("Error decondig response", "Error", err)
		return err
	}

	return nil
}

func req(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request with context", "Error", err)
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
		slog.Error("Error while doing a request", "Error", err)
		return nil, err
	}

	return resp, nil
}

func makeRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request with context", "Error", err)
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

	return req, nil
}

package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"manga-downloader/services"
)

var logger = services.NewLogger()

func ReqImg(ctx context.Context, url string) ([]byte, error) {
	resp, err := req(ctx, url)
	if err != nil {
		log.Println("Error getting response")
		return nil, err
	}
	defer resp.Body.Close()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading image")
		return nil, err
	}

	return img, nil
}

func ReqJSON(ctx context.Context, url string, data any) error {
	resp, err := req(ctx, url)
	if err != nil {
		log.Println("Error getting response")
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		log.Println("Error decondig response")
		return err
	}

	return nil
}

func req(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request with context")
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
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func makeRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request with context")
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
	req.Header = header

	return req, nil
}

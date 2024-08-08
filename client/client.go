package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"

	"mangalib-downloader/logger"
)

type MangaLibClient struct {
	client *http.Client
	Logger *logger.Logger

	Downloaded   chan struct{}
	DownloadPath string

	Page   int
	Query  string
	Branch int
}

func NewClient() *MangaLibClient {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	return &MangaLibClient{
		Logger: logger.NewLogger(),
		client: client,

		Downloaded:   make(chan struct{}, 1),
		DownloadPath: DefaultDownloadPath(),

		Page:  1,
		Query: "",
	}
}

func (c *MangaLibClient) Req(ctx context.Context, url string) (*http.Response, error) {
	req, err := makeRequest(ctx, url)
	if err != nil {
		log.Println("Error creating request with context")
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("Error getting response")
		return nil, err
	}

	return resp, nil
}

func (c *MangaLibClient) ReqImg(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.Req(ctx, url)
	if err != nil {
		log.Println("Error getting response")
		return nil, err
	}
	defer resp.Body.Close()

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.WriteLog(err.Error())
	}

	return img, nil
}

func (c *MangaLibClient) ReqAndDecode(ctx context.Context, url string, data any) error {
	resp, err := c.Req(ctx, url)
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

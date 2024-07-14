package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"mangalib-downloader/logger"
)

var Logger = logger.Logger{}

type MangaLibClient struct {
	client *http.Client
	header http.Header

	Page   int
	Query  string
	Branch int
}

func NewClient() *MangaLibClient {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return &MangaLibClient{
		client: http.DefaultClient,
		header: header,
		Page:   1,
		Query:  "",
	}
}

func (c *MangaLibClient) Req(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request with context")
		return nil, err
	}

	req.Header = c.header

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("Error getting response")
		return nil, err
	}

	return resp, nil
}

func (c *MangaLibClient) ReqAndDecode(ctx context.Context, url string, data any) error {
	resp, err := c.Req(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		fmt.Println("Error decondig response")
		return err
	}

	return nil
}

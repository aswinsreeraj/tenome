package crawler

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"tenome/internal/model"
)

type Crawler interface {
	Crawl(ctx context.Context, url string) (model.Page, error)
}

type HTTPCrawler struct {
	client *http.Client
}

func New(client *http.Client) *HTTPCrawler {
	return &HTTPCrawler{client: client}
}

func (c *HTTPCrawler) Crawl(ctx context.Context, url string) (model.Page, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return model.Page{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model.Page{}, fmt.Errorf("Not Found")
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return model.Page{}, err
	}

	return model.Page{}, nil

}

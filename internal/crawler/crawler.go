package crawler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"tenome/internal/model"

	"golang.org/x/net/html"
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
	if err != nil {
		return model.Page{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return model.Page{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model.Page{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return model.Page{}, err
	}

	title := extractTitle(doc)
	content := extractText(doc)

	return model.Page{URL: url, Title: title, Content: content}, nil

}

func extractTitle(doc *html.Node) string {

	var walk func(node *html.Node)
	var title string

	walk = func(node *html.Node) {
		if title != "" {
			return
		}

		if node.Type == html.ElementNode && node.Data == "title" {
			if node.FirstChild != nil {
				title = node.FirstChild.Data
			}
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}

	}

	walk(doc)

	return title
}

func extractText(doc *html.Node) string {
	var builder strings.Builder

	var walk func(node *html.Node)

	walk = func(node *html.Node) {

		if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "style") {
			return
		}

		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				builder.WriteString(node.Data)
				builder.WriteByte(' ')
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}

	}

	walk(doc)

	return builder.String()

}

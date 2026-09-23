package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx,"GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error:\n\tAn error ocurred forming the Request:\n\t%v\n", err)
	}
	
	req.Header.Set("User-Agent", "You know too much")
	cli := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error:\n\tAn error ocurred doing the request:\n\t%v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Error:\n\tBad http response. Code: %d\n", res.StatusCode)
	}
	
	
	var container RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("Error:\n\tAn error ocurred decoding the response:\n\t%v\n", err)
	}

	cleanResponse := cleanEscapedEntities(container)

	return &cleanResponse, nil
}

func cleanEscapedEntities(Rssfeed RSSFeed) RSSFeed {
	fmt.Printf("\n\nFetched!\n\n")
	Rssfeed.Channel.Title = html.UnescapeString(Rssfeed.Channel.Title)
	Rssfeed.Channel.Description = html.UnescapeString(Rssfeed.Channel.Description)
	for i := range Rssfeed.Channel.Item {
		Rssfeed.Channel.Item[i].Title = html.UnescapeString(Rssfeed.Channel.Item[i].Title)
		Rssfeed.Channel.Item[i].Description = html.UnescapeString(Rssfeed.Channel.Item[i].Description)
	}
	return Rssfeed
}
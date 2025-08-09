package main

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, errors.New("FAILED TO CREATE REQUEST: " + err.Error())
	}

	req.Header.Set("User-Agent", "gator")

	res, err2 := client.Do(req)
	if err2 != nil {
		return nil, errors.New("FAILED TO FETCH FEED: " + err2.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("FAILED TO FETCH FEED: " + res.Status)
	}

	body, err3 := io.ReadAll(res.Body)
	if err3 != nil {
		return nil, errors.New("FAILED TO READ FEED BODY: " + err3.Error())
	}

	var rssFeed RSSFeed
	if err4 := xml.Unmarshal(body, &rssFeed); err4 != nil {
		return nil, errors.New("FAILED TO UNMARSHAL FEED: " + err4.Error())
	}

	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)

	for _, item := range rssFeed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
	}

	return &rssFeed, nil
}

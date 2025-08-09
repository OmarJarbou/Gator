package main

import (
	"context"
	"errors"
	"fmt"
)

func scrapeFeeds(ctx context.Context, state *state) error {
	nextFeed, err := state.DBQueries.GetNextFeedToFetch(ctx)
	if err != nil {
		return errors.New("FAILED TO GET NEXT FEED TO FETCH: " + err.Error())
	}

	rssFeed, err2 := fetchFeed(context.Background(), nextFeed.Url)
	if err2 != nil {
		return err2
	}

	err3 := state.DBQueries.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err3 != nil {
		return errors.New("FAILED TO MARK FEED AS FETCHED: " + err3.Error())
	}

	fmt.Println("Feed", nextFeed.Name, "fetched successfully!")
	fmt.Println("Items:")
	for _, item := range rssFeed.Channel.Item {
		fmt.Println("-", item.Title)
	}

	return nil
}

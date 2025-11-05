package main

import (
	"context"
	"errors"
	"fmt"

	"database/sql"
	"time"

	"github.com/OmarJarbou/Gator/internal/database"
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
		post := database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: func() sql.NullTime {
				t, err := time.Parse(time.RFC1123Z, item.PubDate)
				if err != nil {
					return sql.NullTime{Valid: false}
				}
				return sql.NullTime{Time: t, Valid: true}
			}(),
			FeedID: nextFeed.ID,
		}
		_, err4 := state.DBQueries.CreatePost(context.Background(), post)
		if err4 != nil {
			return errors.New("FAILED TO CREATE POST: " + err4.Error())
		}
		fmt.Println("Post created successfully!")
	}

	return nil
}

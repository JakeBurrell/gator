package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JakeBurrell/gator/internal/database"
	"github.com/JakeBurrell/gator/internal/rss"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get feed to fetch", err)
		return
	}
	log.Println("Found a feed to fetch")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s: %v", feed.Name, err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

}

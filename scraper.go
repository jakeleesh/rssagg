package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jakeleesh/rssagg/internal/database"
)

// Scrapper is a long running job, running background as Server runs
// Inputs:
// A connection to database
// Concurrency units: How many goroutines want to do the scraping
// How much time we want in between each request to go scrape a new RSSFeed
// Shouldn't return anything because going to be a long running job
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	// Scraper running in background of server, important have good logging, tells us what's going on
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	// Make request on interval
	// Responds with a ticker
	ticker := time.NewTicker(timeBetweenRequest)

	// Use for loop to execute every time new value comes across ticker's channel
	// ticker has field C, which is a channel where value sent across the channel
	// Passing in empty initializes and middle is so that it executes immediately
	// If did for range ticker.C, will wait
	for ; ; <-ticker.C {
		// context.Background() is global context
		// Use if don't have access to scoped context like for individual http requests
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds:", err)
			// Continue because function should always be running as server operates
			continue
		}

		// Fetches each feed individually at the same time
		// Need synchronization mechanism: Use WaitGroup
		wg := &sync.WaitGroup{}
		// Iterating over all feeds we want to fetch on individual goroutines
		for _, feed := range feeds {
			// The way WaitGroup works:
			// Anytime want spawn new goroutine within context of WaitGroup, add number to it
			// Iterating over all feeds on the same goroutine startScraping function
			// Adding 1 to WaitGroup for every feed
			// Had concurrency of 30, adding 30 to WaitGroup
			// Spawning seperate goroutines
			// End of loop, waiting on WaitGroup for 30 distinct calls to wg.Done()
			// wg.Done() decrements counter by 1, adding 1 everytime iterate over slice
			// Calling done when done scraping feed
			// ALlows us to scrape feed same time 30 times
			// Spawn 30 different gorountines to scrape 30 different RSSFeed
			wg.Add(1)

			// Spawn new goroutine, pass WaitGroup in
			go scrapeFeed(db, wg, feed)
		}
		// When all done, will execute
		// Before done, will be blocking
		// Don't want to continue next iteration of loop until sure scraped all feeds
		wg.Wait()
	}
}

// Pointer to WaitGroup
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// Decrements counter by 1
	// Deferring so will always be called at end of function
	defer wg.Done()

	// Mark that we're fetching this feed
	// Returns updated feed, can ignore
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		// Not returning anything from function, calling on new goroutine so nothign to return
		// Just log there was an issue
		log.Println("Error marking feed as fetched:", err)
		return
	}

	// Scrape Feed
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		// NullString has string itself and whether it's valid
		// If item description is blank, set the value to null in database
		// Otherwise create valid description
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Need to parse
		// RFC1123Z is layout
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}

		// Don't care about new post
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Scrapped feeds and pulled in posts
			// "duplicate key violates unique constraints" makes sense, didn't want to store duplicate posts in database
			// Try to recreate post, fails because already have posts in database
			// String detection, don't log this, isn't an error, expected behavior
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			// Log error if not "duplicate key" error
			log.Println("failed to create post", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

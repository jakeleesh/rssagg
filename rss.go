package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// Keys for RSS entries in https://www.wagslane.dev/ blog
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Parse
func urlToFeed(url string) (RSSFeed, error) {
	// Need HTTP Client
	// Create using http library
	httpClient := http.Client{
		// Set Timeout to 10s, more than 10s to fetch, don't want, probably broken
		Timeout: 10 * time.Second,
	}

	// Use Client to make GET request to URL of feed
	// Return http response
	resp, err := httpClient.Get(url)
	if err != nil {
		// Return empty structs
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	// Get all data from response body
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}

	// Want to read into RSSFeed
	// Similar to dealing with JSON
	// Pointer to where we want to unmarshal the data into that location in memory
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	// Can just return populated RSSFeed
	return rssFeed, nil
}

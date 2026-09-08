package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/turkerkiv/gator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, val := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(val.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(val.Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	feedToNext, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID:        feedToNext.ID,
	}
	if err := s.db.MarkFeedFetched(context.Background(), params); err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feedToNext.Url)
	if err != nil {
		return err
	}

	for _, val := range feed.Channel.Item {
		if val.Link == "" {
			continue
		}

		t, err := time.Parse(time.RFC1123, val.PubDate)
		timeParam := sql.NullTime{}
		if err == nil {
			timeParam = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		descParam := sql.NullString{}
		if val.Description != "" {
			descParam = sql.NullString{
				String: val.Description,
				Valid:  true,
			}
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      feedToNext.ID,
			Title:       val.Title,
			Url:         val.Link,
			Description: descParam,
			PublishedAt: timeParam,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			continue
		}
		fmt.Printf("collected - url: %s - feed: %s\n", feedToNext.Url, val.Title)
	}

	return nil
}

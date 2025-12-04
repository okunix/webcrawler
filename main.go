package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	ctx := context.TODO()

	// initializing fetcher
	fetcher := NewDefaultFetcher(UserAgent, Timeout)

	fetched := make(map[string]bool)
	linksCh := make(chan []string)
	unseenLinkCh := make(chan string)

	// loading urls from os.Args
	go func() { linksCh <- os.Args[1:] }()

	// starting 20 workers
	for range 20 {
		go func() {
			for link := range unseenLinkCh {
				fetchedLinks, code, err := fetcher.Fetch(ctx, link)
				if err != nil {
					slog.Error("error while fetching",
						"err", err.Error(),
						"code", code,
						"url", link,
					)
					return
				}
				slog.Info("fetch", "url", link, "code", code)
				go func() { linksCh <- fetchedLinks }()
			}
		}()
	}

	// checking if we have seen link
	// if not add it to unseenLink channel and mark as seen
	for links := range linksCh {
		for _, link := range links {
			if fetched[link] {
				continue
			}
			fetched[link] = true
			unseenLinkCh <- link
		}
	}
}

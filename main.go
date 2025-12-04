package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	// initializing fetcher
	timeoutSec, err := strconv.Atoi(Timeout)
	if err != nil {
		panic(err)
	}
	timeout := time.Duration(timeoutSec) * time.Second
	ctx := context.TODO()
	fetcher := NewDefaultFetcher(UserAgent, timeout)

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

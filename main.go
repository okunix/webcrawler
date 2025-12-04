package main

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"
)

func main() {
	initURLs := os.Args[1:]
	// getting website domains to follow
	baseURLs := []string{}
	for _, v := range initURLs {
		url, err := BaseURL(v)
		if err != nil {
			panic(err)
		}
		baseURLs = append(baseURLs, url)
	}

	fetched := make(map[string]bool)

	timeoutSec, err := strconv.Atoi(Timeout)
	if err != nil {
		panic(err)
	}
	timeout := time.Duration(timeoutSec) * time.Second

	ctx := context.TODO()
	fetcher := NewDefaultFetcher(UserAgent, timeout)

	var n int
	linksCh := make(chan []string)
	go func() { linksCh <- initURLs }()

	sem := make(chan struct{}, 20)
	var wg sync.WaitGroup
	n++
	for ; n > 0; n-- {
		links := <-linksCh
		for _, link := range links {
			baseURL, _ := BaseURL(link)
			if !slices.Contains(baseURLs, baseURL) {
				continue
			}

			if fetched[link] {
				continue
			}
			fetched[link] = true
			sem <- struct{}{}
			wg.Add(1)
			n++
			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()
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
				linksCh <- fetchedLinks
			}()
		}
	}
	wg.Wait()
}

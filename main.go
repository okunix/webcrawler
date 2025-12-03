package main

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"time"
)

func main() {
	initURLs := os.Args[1:]
	fetched := make(map[string]bool)

	timeoutSec, err := strconv.Atoi(Timeout)
	if err != nil {
		panic(err)
	}

	timeout := time.Duration(timeoutSec) * time.Second

	ctx := context.TODO()
	fetcher := NewDefaultFetcher(UserAgent, timeout)
	links := NewQueue[string]()
	links.Enqueue(initURLs...)

	// getting website domains to follow
	baseURLs := []string{}
	for _, v := range initURLs {
		url, err := BaseURL(v)
		if err != nil {
			panic(err)
		}
		baseURLs = append(baseURLs, url)
	}

	for !links.Empty() {
		url, err := links.Dequeue()
		if err != nil {
			continue
		}

		baseURL, _ := BaseURL(url)
		if !slices.Contains(baseURLs, baseURL) {
			continue
		}

		if fetched[url] {
			continue
		}
		fetched[url] = true

		fetchedLinks, code, err := fetcher.Fetch(ctx, url)
		if err != nil {
			slog.Error("error while fetching",
				"err", err.Error(),
				"code", code,
				"url", url,
			)
			continue
		}
		slog.Info("fetch", "url", url, "code", code)
		links.Enqueue(fetchedLinks...)
	}
}

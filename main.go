package main

import (
	"context"
	"log/slog"
	"os"
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

	linksCh := make(chan []string)
	go func() { linksCh <- initURLs }()

	// getting website domains to follow
	//baseURLs := []string{}
	//for _, v := range initURLs {
	//url, err := BaseURL(v)
	//if err != nil {
	//panic(err)
	//}
	//baseURLs = append(baseURLs, url)
	//}
	for links := range linksCh {
		for _, link := range links {
			//baseURL, _ := BaseURL(link)
			//if !slices.Contains(baseURLs, baseURL) {
			//continue
			//}

			if fetched[link] {
				continue
			}
			fetched[link] = true
			go func() {
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
}

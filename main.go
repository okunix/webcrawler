package main

import (
	"context"
	"log/slog"
	"os"
	"slices"

	"github.com/okunix/webcrawler/config"
	"github.com/okunix/webcrawler/fetcher"
	linkPkg "github.com/okunix/webcrawler/link"
)

func main() {
	ctx := context.TODO()

	domains := []string{}
	for _, v := range os.Args[1:] {
		domain, err := linkPkg.BaseURL(v)
		if err != nil {
			panic(err)
		}
		domains = append(domains, domain)
	}

	// initializing fetcher
	fetcher := fetcher.NewDefaultFetcher(config.UserAgent, config.Timeout)

	seen := make(map[string]bool)
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
	var n int
	for links := range linksCh {
		for _, link := range links {
			domain, err := linkPkg.BaseURL(link)
			if err != nil {
				continue
			}
			if seen[link] || !slices.Contains(domains, domain) {
				continue
			}
			seen[link] = true
			n++
			unseenLinkCh <- link
		}
		n--
		if n < 0 {
			break
		}
	}
}

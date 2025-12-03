package main

import (
	"errors"
	"regexp"
)

var (
	hrefRegex         = regexp.MustCompile(`<a\s+href=["'](https?://[^"'\s]+|(?:/[^"'\s]+)+)["']>`)
	baseURLRegex      = regexp.MustCompile(`^(https?://[^/\\<>#\s]+)`)
	relativePathRegex = regexp.MustCompile(`^(?:/[^/\\\s]+)+/?(?:\?.+)?$`)
	httpRegex         = regexp.MustCompile(`^https?://\S+`)
)

func ExtractLinks(s, baseURL string) []string {
	links := ExtractRawLinks(s)
	return Normalize(baseURL, links)
}

func ExtractRawLinks(s string) []string {
	links := []string{}
	hrefs := hrefRegex.FindAllStringSubmatch(s, -1)
	if hrefs == nil {
		return links
	}
	for _, href := range hrefs {
		links = append(links, href[1])
	}
	return links
}

func Normalize(baseURL string, links []string) []string {
	normalized := []string{}
	for _, v := range links {
		if httpRegex.MatchString(v) {
			normalized = append(normalized, v)
		} else if relativePathRegex.MatchString(v) {
			normalized = append(normalized, baseURL+v)
		}
	}
	return normalized
}

func BaseURL(link string) (string, error) {
	matches := baseURLRegex.FindStringSubmatch(link)
	if len(matches) < 2 {
		return "", errors.New("base of the url not detected")
	}
	return matches[1], nil
}

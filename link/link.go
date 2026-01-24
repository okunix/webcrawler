package link

import (
	"errors"
	"regexp"
)

var (
	hrefRegex = regexp.MustCompile(
		`<a\s+href=["']?(https?://[^"'>\s]+|(?:/[^"'>\s]+)+)["']?>`,
	)
	baseURLRegex      = regexp.MustCompile(`^(https?://[^/\\<>#\s]+)`)
	domainRegex       = regexp.MustCompile(`^https?://([^/\\<>#\s]+)`)
	relativePathRegex = regexp.MustCompile(`^(?:/[^/\\\s]+)+/?(?:\?.+)?$`)
	httpRegex         = regexp.MustCompile(`^https?://\S+`)
	htmlIdRefRegex    = regexp.MustCompile(`#[A-Za-z][A-Za-z0-9\-_:.]*$`)
)

func Extract(body, baseURL string) []string {
	links := ExtractHrefs(body)
	return Normalize(baseURL, links)
}

func ExtractHrefs(body string) []string {
	links := []string{}
	hrefs := hrefRegex.FindAllStringSubmatch(body, -1)
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
		var link string
		if httpRegex.MatchString(v) {
			link = v
		} else if relativePathRegex.MatchString(v) {
			link = baseURL + v
		} else {
			continue
		}
		cleanLink := htmlIdRefRegex.ReplaceAllString(link, "")
		normalized = append(normalized, cleanLink)
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

func Domain(link string) (string, error) {
	matches := domainRegex.FindStringSubmatch(link)
	if len(matches) < 2 {
		return "", errors.New("domain not detected")
	}
	return matches[1], nil
}

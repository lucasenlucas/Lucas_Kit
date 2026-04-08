package crawl

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// CrawlResult contains links, scripts, forms.
type CrawlResult struct {
	Links   []string `json:"links,omitempty"`
	Scripts []string `json:"scripts,omitempty"`
	Forms   []string `json:"forms,omitempty"`
}

// GetLinks extracts href attributes from <a> tags.
func GetLinks(body io.Reader) []string {
	var links []string
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// GetScripts extracts src attributes from <script> tags.
func GetScripts(body io.Reader) []string {
	var scripts []string
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return scripts
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "script" {
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						scripts = append(scripts, attr.Val)
					}
				}
			}
		}
	}
}

// GetForms extracts action and method from <form> tags.
func GetForms(body io.Reader) []string {
	var forms []string
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return forms
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "form" {
				action := ""
				method := "GET"
				for _, attr := range token.Attr {
					if attr.Key == "action" {
						action = attr.Val
					}
					if attr.Key == "method" {
						method = strings.ToUpper(attr.Val)
					}
				}
				forms = append(forms, fmt.Sprintf("%s [%s]", action, method))
			}
		}
	}
}

// GetFavicon attempts to find the favicon URL.
func GetFavicon(domain string, body io.Reader) string {
	favURL := fmt.Sprintf("https://%s/favicon.ico", domain)
	resp, err := http.Head(favURL)
	if err == nil && resp.StatusCode == 200 {
		return favURL
	}

	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return ""
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "link" {
				isIcon := false
				href := ""
				for _, attr := range token.Attr {
					if attr.Key == "rel" && (strings.Contains(attr.Val, "icon") || strings.Contains(attr.Val, "shortcut")) {
						isIcon = true
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if isIcon && href != "" {
					return href
				}
			}
		}
	}
}

// GetTitle extracts the page title.
func GetTitle(body io.Reader) string {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return ""
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					return tokenizer.Token().Data
				}
			}
		}
	}
}

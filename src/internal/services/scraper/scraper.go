package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"ecommerce-price-tracker/internal/models"
)

type Scraper interface {
	// Identifier is the registry key, matched against URL hosts.
	Identifier() string
	DefaultSelectors() Selectors
	// GetProductInfo returns the full product details for the given URL.
	GetProductInfo(url string) (models.ProductInfo, error)
	// GetPrice returns only the current price for the given URL.
	GetPrice(url string) (float64, error)
	StripURL(url string) string
}

var registry = map[string]Scraper{}

func Register(s Scraper) {
	registry[s.Identifier()] = s
}

func Get(id string) (Scraper, error) {
	s, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("no scraper registered with identifier %q", id)
	}
	return s, nil
}

func For(rawURL string) (Scraper, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid product url %q: %w", rawURL, err)
	}
	host := u.Hostname()
	if host == "" {
		host = rawURL // tolerate URLs passed without a scheme
	}
	for id, s := range registry {
		if strings.Contains(host, id) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("no scraper supports url %q", rawURL)
}

type BaseScraper struct{}

func (BaseScraper) StripURL(url string) string {
	if idx := strings.Index(url, "?"); idx != -1 {
		return url[:idx]
	}
	return url
}

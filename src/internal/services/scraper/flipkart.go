package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
)

const FlipkartID = "flipkart"

func init() {
	Register(NewFlipkartScraper())
}

type FlipkartScraper struct {
	BaseScraper
}

func NewFlipkartScraper() *FlipkartScraper {
	return &FlipkartScraper{}
}

func (s *FlipkartScraper) Identifier() string {
	return FlipkartID
}

func (s *FlipkartScraper) DefaultSelectors() Selectors {
	return Selectors{
		Price:     "div.v1zwn20:nth-child(1)",
		Name:      "h1.v1zwn21m",
		Image:     "div._7dzyg24w:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > picture:nth-child(1) > img:nth-child(3)",
		ImageAttr: "src",
	}
}

func (s *FlipkartScraper) GetProductInfo(url string) (models.ProductInfo, error) {
	url = s.StripURL(url)
	sels := loadSelectors(s.Identifier(), s.DefaultSelectors())
	doc, err := fetchDocument(url)
	if err != nil {
		return models.ProductInfo{}, err
	}

	priceText := ""
	doc.Find(sels.Price).Each(func(i int, sel *goquery.Selection) {
		priceText = sel.Text()
	})
	if priceText == "" {
		return models.ProductInfo{}, errors.New("error: unable to extract price")
	}
	price, err := parsePrice(priceText)
	if err != nil {
		return models.ProductInfo{}, err
	}

	name := ""
	doc.Find(sels.Name).Each(func(i int, sel *goquery.Selection) {
		name = sel.Text()
	})
	if name == "" {
		// Flipkart's obfuscated h1 class hash drifts between deploys, breaking the
		// configured selector. The product title is the page's sole <h1>, so fall
		// back to that before giving up.
		name = strings.TrimSpace(doc.Find("h1").First().Text())
	}
	if name == "" {
		return models.ProductInfo{}, errors.New("error: unable to extract ProductInfo name")
	}

	imgUrl := ""
	doc.Find(sels.Image).Each(func(i int, sel *goquery.Selection) {
		src, ok := sel.Attr(sels.ImageAttr)
		if !ok {
			src = "" // TODO: add placeholder image
		}
		imgUrl = src
	})

	return models.ProductInfo{
		Name:    name,
		Price:   price,
		ImgLink: imgUrl,
		Url:     url,
	}, nil
}

func (s *FlipkartScraper) GetPrice(url string) (float64, error) {
	url = s.StripURL(url)
	sels := loadSelectors(s.Identifier(), s.DefaultSelectors())
	doc, err := fetchDocument(url)
	if err != nil {
		return 0, err
	}

	priceText := ""
	doc.Find(sels.Price).Each(func(i int, sel *goquery.Selection) {
		priceText = sel.Text()
	})
	if priceText == "" {
		return 0, errors.New("error: unable to extract price")
	}
	return parsePrice(priceText)
}

var priceCleanup = regexp.MustCompile(`(?i)rs|₹|[.,]`)

func parsePrice(priceText string) (float64, error) {
	priceText = priceCleanup.ReplaceAllString(priceText, "")
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing price: %w", err)
	}
	return price, nil
}

func fetchDocument(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}
	return doc, nil
}

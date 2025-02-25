package services

import (
	"ecommerce-price-tracker/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetNameAndPrice(url string) (models.Product, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Product{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Product{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Product{}, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Product{}, fmt.Errorf("error parsing HTML: %w", err)
	}

	var product models.Product

	priceText := ""
	doc.Find(".VU-ZEz").Each(func(i int, s *goquery.Selection) {
		priceText = s.Text()
	})

	if priceText == "" {
		return models.Product{}, errors.New("error: unable to extract price")
	}

	priceText = strings.ReplaceAll(priceText, ",", "")
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return models.Product{}, fmt.Errorf("error parsing price: %w", err)
	}

	name := ""
	doc.Find(".CxhGGd").Each(func(i int, s *goquery.Selection) {
		name = s.Text()
	})

	if name == "" {
		return models.Product{}, errors.New("error: unable to extract product name")
	}

	product.Name = name
	product.Price = price

	return product, nil
}

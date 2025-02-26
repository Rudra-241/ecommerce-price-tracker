package services

import (
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func GetProductInfo(url string) (models.ProductInfo, error) {
	url = StripURL(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return models.ProductInfo{}, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error parsing HTML: %w", err)
	}

	var ProductInfo models.ProductInfo

	priceText := ""
	doc.Find(".CxhGGd").Each(func(i int, s *goquery.Selection) {
		priceText = s.Text()
	})

	if priceText == "" {
		return models.ProductInfo{}, errors.New("error: unable to extract price")
	}

	re := regexp.MustCompile(`(?i)rs|₹|[.,]`)
	priceText = re.ReplaceAllString(priceText, "")
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error parsing price: %w", err)
	}

	name := ""
	doc.Find(".VU-ZEz").Each(func(i int, s *goquery.Selection) {
		name = s.Text()
	})
	if name == "" {
		return models.ProductInfo{}, errors.New("error: unable to extract ProductInfo name")
	}

	imgUrl := ""
	var ok bool
	doc.Find(".jLEJ7H").Each(func(i int, s *goquery.Selection) {
		imgUrl, ok = s.Attr("src")
		if !ok {
			imgUrl = "" // TODO: add placeholder image

		}
	})

	ProductInfo.Name = name
	ProductInfo.Price = price
	ProductInfo.ImgLink = imgUrl
	ProductInfo.Url = url
	return ProductInfo, nil
}

func GetPrice(url string) (float64, error) {
	url = StripURL(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error parsing HTML: %w", err)
	}

	priceText := ""
	doc.Find(".CxhGGd").Each(func(i int, s *goquery.Selection) {
		priceText = s.Text()
	})

	if priceText == "" {
		return 0, errors.New("error: unable to extract price")
	}

	re := regexp.MustCompile(`(?i)rs|₹|[.,]`)
	priceText = re.ReplaceAllString(priceText, "")
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing price: %w", err)
	}
	return price, nil
}

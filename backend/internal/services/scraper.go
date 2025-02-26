package services

import (
	"ecommerce-price-tracker/internal/models"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func GetNameAndPrice(url string) (models.ProductInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ProductInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0")
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
	fmt.Print(priceText)
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
			imgUrl = ""
		}
	})

	ProductInfo.Name = name
	ProductInfo.Price = price
	ProductInfo.ImgLink = imgUrl
	fmt.Println(ProductInfo)
	return ProductInfo, nil
}

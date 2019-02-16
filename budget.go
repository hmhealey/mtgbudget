package mtgtools

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	searchRoot = "https://shop.tcgplayer.com/productcatalog/product/show?ProductType=All&IsProductNameExact=false&ProductName="
)

var (
	ErrNoPrice = errors.New("No price available")
)

func GetMarketPrice(cardName string) (string, int, error) {
	c := colly.NewCollector()

	prices := make(map[string]int)

	c.OnHTML(".product__card", func(e *colly.HTMLElement) {
		if cardName != e.ChildText(".product__name") {
			return
		}

		set := e.ChildText(".product__group")
		if !strings.HasSuffix(set, " (Magic)") {
			return
		}

		set = set[:len(set)-8]

		if set == "World Championship Decks" {
			return
		}

		rawPrice := e.ChildText(".product__market-price dd")
		if rawPrice == "Unavailable" {
			return
		}

		rawPriceSplit := strings.SplitN(rawPrice, ".", 2)
		if len(rawPriceSplit) != 2 {
			log.Printf("Entry for %s with set %s has invalid price %s", cardName, set, rawPrice)
			return
		}

		dollars, err := strconv.ParseInt(rawPriceSplit[0][1:], 10, 32)
		if err != nil {
			log.Printf("Entry for %s with set %s has invalid dollar amount %s from %s", cardName, set, rawPriceSplit[0][1:], rawPrice)
			return
		}

		cents, err := strconv.ParseInt(rawPriceSplit[1], 10, 32)
		if err != nil {
			log.Printf("Entry for %s with set %s has invalid cent amount %s from %s", cardName, set, rawPriceSplit[1], rawPrice)
			return
		}

		prices[set] = int(dollars)*100 + int(cents)
	})

	err := c.Visit(fmt.Sprint(searchRoot, url.QueryEscape(cardName)))
	if err != nil {
		log.Fatalf("Failed to visit web page for %s: %v", cardName, err)
	}

	lowPrice := 0
	lowSet := ""

	for set, price := range prices {
		if lowSet == "" || price < lowPrice {
			lowPrice = price
			lowSet = set
		}
	}

	if lowSet == "" {
		return "", 0, ErrNoPrice
	}

	return lowSet, lowPrice, nil
}

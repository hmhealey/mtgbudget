package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const NO_PRICE = 10000000
const SEARCH_ROOT = "https://shop.tcgplayer.com/productcatalog/product/show?ProductType=All&IsProductNameExact=false&ProductName="

func mainBudget(cardNames []string) {
	var total int

	for _, cardName := range cardNames {
		set, price := getMarketPrice(cardName)

		if set == "" {
			log.Printf("No prices for %s", cardName)
		}

		fmt.Printf("%s (%s): %.2f\n", cardName, set, float32(price) / 100)

		if price != NO_PRICE {
			total += price
		}
	}

	fmt.Println("=======================================")
	fmt.Printf("Total: %.2f\n", float32(total) / 100)
}

func getMarketPrice(cardName string) (string, int) {
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

		prices[set] = int(dollars) * 100 + int(cents)
	})

	err := c.Visit(fmt.Sprint(SEARCH_ROOT, url.QueryEscape(cardName)))
	if err != nil {
		log.Fatalf("Failed to visit web page for %s: %v", cardName, err)
	}

	lowPrice := NO_PRICE
	lowSet := ""

	for set, price := range prices {
		if price < lowPrice {
			lowPrice = price
			lowSet = set
		}
	}

	return lowSet, lowPrice
}
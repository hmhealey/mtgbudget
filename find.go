package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

const FIND_ROOT = "https://deckbox.org/mtg/"

type cardCount struct {
	Set   string
	Count int
}

func mainFindCards(cardNames []string) {
	var notFound []string

	for _, cardName := range cardNames {
		sets, counts := findCardCount(cardName)

		if len(counts) > 0 {
			fmt.Printf("%s\n", cardName)

			for i, count := range counts {
				set := sets[i]

				fmt.Printf("  %s: %d %s found\n", set, count, copies(count))
			}
		} else {
			notFound = append(notFound, cardName)
		}
	}

	if len(notFound) > 0 {
		fmt.Println("=======================================")

		for _, cardName := range notFound {
			fmt.Printf("%s: No copies found\n", cardName)
		}
	}
}

func findCardCount(cardName string) ([]string, []int) {
	c := colly.NewCollector()

	c.SetCookies("https://deckbox.org", []*http.Cookie{
		{
			Name:   "auth_token",
			Value:  DECKBOX_AUTH_TOKEN,
			Domain: "deckbox.org",
		},
	})

	var sets []string
	var counts []int

	c.OnHTML(".left_card_col .warning", func(e *colly.HTMLElement) {
		log.Print("Not logged in")
	})

	c.OnHTML("#in_your_collection tr", func(e *colly.HTMLElement) {
		set := e.ChildAttr(".mtg_edition_container img", "data-title")
		set = regexp.MustCompile(` \(Card #\d+\)$`).ReplaceAllString(set, "")

		if set == "" {
			log.Printf("Failed to get set for card %s", cardName)
			return
		}

		if e.ChildAttr(".sprite.s_colors", "data-title") == "Foil" {
			set += " (Foil)"
		}

		rawCount := e.ChildText(".inventory_count")
		count, err := strconv.ParseInt(rawCount, 10, 32)
		if err != nil {
			log.Printf("Failed to parse count %s for card %s in set %s", rawCount, cardName, set)
			return
		}

		sets = append(sets, set)
		counts = append(counts, int(count))
	})

	err := c.Visit(fmt.Sprint(FIND_ROOT, url.QueryEscape(cardName)))
	if err != nil {
		log.Fatalf("Failed to visit web page for %s: %v", cardName, err)
	}

	return sets, counts
}

func copies(count int) string {
	if count == 1 {
		return "copy"
	} else {
		return "copies"
	}
}
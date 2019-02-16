package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hmhealey/mtgtools"
)

func mainBudget(args []string) {
	cardNames := readCardNames(os.Stdin)

	var total int

	for _, cardName := range cardNames {
		set, price, err := mtgtools.GetMarketPrice(cardName)

		if err == mtgtools.ErrNoPrice {
			log.Printf("%s: No price available\n", cardName)
			continue
		}

		fmt.Printf("%s (%s): %.2f\n", cardName, set, float32(price)/100)

		total += price
	}

	fmt.Println("=======================================")
	fmt.Printf("Total: %.2f\n", float32(total)/100)
}

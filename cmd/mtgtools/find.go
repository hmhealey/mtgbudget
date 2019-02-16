package main

import (
	"fmt"
	"os"

	"github.com/hmhealey/mtgtools"
)

func mainFind(args []string) {
	cardNames := readCardNames(os.Stdin)

	var notFound []string

	for _, cardName := range cardNames {
		sets, counts := mtgtools.FindCardCount(cardName)

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

func copies(count int) string {
	if count == 1 {
		return "copy"
	}

	return "copies"
}

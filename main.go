package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var IGNORED_CARDS = map[string]bool{
	"Plains": true,
	"Island": true,
	"Swamp": true,
	"Mountain": true,
	"Forest": true,
}

func parseCardNames(deckList string) []string {
	var cardNames []string

	for _, line := range strings.Split(deckList, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		line = regexp.MustCompile(`^\d+x `).ReplaceAllString(line, "")
		line = regexp.MustCompile(` \*CMDR\*`).ReplaceAllString(line, "")
		line = regexp.MustCompile(` #\w+`).ReplaceAllString(line, "")

		cardName := strings.TrimSpace(line)

		if _, ok := IGNORED_CARDS[cardName]; ok {
			continue
		}

		cardNames = append(cardNames, cardName)
	}

	return cardNames
}

func main() {
	if len(os.Args) == 1 {
		log.Print("No mode specified")
		return
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Print("Failed to read deck list from stdin")
		return
	}

	deckList := string(b)
	if deckList == "" {
		log.Print("No deck list provided")
		return
	}

	cardNames := parseCardNames(deckList)

	mode := os.Args[1]
	switch mode {
	case "budget":
		mainBudget(cardNames)
	case "find":
		mainFindCards(cardNames)
	}
}
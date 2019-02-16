package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var ignoredCards = map[string]bool{
	"Plains":   true,
	"Island":   true,
	"Swamp":    true,
	"Mountain": true,
	"Forest":   true,
}

func readCardNames(r io.Reader) []string {
	scanner := bufio.NewScanner(r)

	var cardNames []string

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		line = regexp.MustCompile(`^\d+x `).ReplaceAllString(line, "")
		line = regexp.MustCompile(` \*CMDR\*`).ReplaceAllString(line, "")
		line = regexp.MustCompile(` #\w+`).ReplaceAllString(line, "")

		cardName := strings.TrimSpace(line)

		if _, ok := ignoredCards[cardName]; ok {
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

	mode := os.Args[1]
	args := os.Args[2:]
	switch mode {
	case "budget":
		mainBudget(args)
	case "find":
		mainFind(args)
	}
}

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Player struct holds the fields scraped from a squad table row
type Player struct {
	Number   string
	Name     string
	Nation   string
	Position string
	DOB      string
}

// Helper function to strip Wikipedia's footnote and annotation markups
func cleanName(s string) string {
	re := regexp.MustCompile(`[\*†#◊]|\(.*?\)|\[.*?\]`)
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}

// Function that scrapes data from given wikipedia team site and returns the collected players
func scrape(link string) []Player {
	// Slice to collect players as we find them
	var players []Player

	// Initialize colly
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// On every wikitable element, check if it's the squad list and extract player rows
	c.OnHTML("table.wikitable", func(table *colly.HTMLElement) {
		header := table.ChildText("tr:first-child")

		// Only process the table that looks like a squad list
		if !strings.Contains(header, "No.") || !strings.Contains(header, "Player") || !strings.Contains(header, "Nat.") {
			return
		}

		// Loop over every row in the matched squad table
		table.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			// Collect the text of each <td> cell in this row
			var cells []string
			row.DOM.Find("td").Each(func(_ int, s *goquery.Selection) {
				cells = append(cells, strings.TrimSpace(s.Text()))
			})

			// Skip header row and section dividers which don't have the full set of columns
			if len(cells) < 9 {
				return
			}

			// Build a Player from the row and append it to our results
			players = append(players, Player{
				Number:   cells[0],
				Name:     cleanName(cells[1]),
				Nation:   cells[2],
				Position: cells[3],
				DOB:      cells[4],
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(link)

	return players
}

// Print out every player in given splice
func printPlayers(players []Player) {
	for _, p := range players {
		fmt.Printf("#%s | %s | %s | %s | %s\n", p.Number, p.Name, p.Nation, p.Position, p.DOB)
	}
}

func main() {
	players := scrape("https://en.wikipedia.org/wiki/2024%E2%80%9325_Arsenal_F.C._season")
	printPlayers(players)
}

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Helper function to strip Wikipedia's footnote and annotation markups
func cleanName(s string) string {
	re := regexp.MustCompile(`[\*†#◊]|\(.*?\)|\[.*?\]`)
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}

// Function that scrapes data from given wikipedia team site
func scrape(link string) {
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

			// Pull out the fields we care about by column position
			number := cells[0]
			player := cleanName(cells[1])
			nation := cells[2]
			position := cells[3]
			dob := cells[4]

			// Print link found
			fmt.Printf("#%s | %s | %s | %s | %s\n", number, player, nation, position, dob)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(link)
}

func main() {
	scrape("https://en.wikipedia.org/wiki/2024%E2%80%9325_Arsenal_F.C._season")
}

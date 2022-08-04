package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Response struct {
	Rankings []Ranking
	Events   []Event
}

type Ranking struct {
	WeightClass string
	Weight      string
	Fighters    []RankedFighter
}

type RankedFighter struct {
	Name        string
	Rank        string
	RankChange  string
	SherdogLink string
}

type Event struct {
	Name           string
	Date           string
	Venue          string
	Location       string
	MainCardFights []Fight
	PrelimFights   []Fight
}

type Fight struct {
	FighterOne  Fighter
	FighterTwo  Fighter
	WeightClass string
}

type Fighter struct {
	Name       string
	Rank       string
	Wins       int
	Losses     int
	Draws      int
	NoContests int
}

func Scraper() Response {
	response := Response{}

	response.Rankings = ScrapeRankings()
	response.Events = ScrapeUpcomingEvents()

	return response
}

func ScrapeRankings() []Ranking {
	fmt.Println("SCRAPING RANKINGS")

	rankings := []Ranking{}

	rankingsLink := "https://en.wikipedia.org/wiki/UFC_Rankings"

	weightClasses := []string{
		"Mens Pound For Pound",
		"Womens Pound For Pound",
		"Heavyweight",
		"Light Heavyweight",
		"Middleweight",
		"Welterweight",
		"Lightweight",
		"Featherweight",
		"Bantamweight",
		"Flyweight",
		"Womens Bantamweight",
		"Womens Flyweight",
		"Womens Strawweight",
	}

	weights := []string{
		"",
		"",
		"265 lbs.",
		"205 lbs.",
		"185 lbs.",
		"170 lbs.",
		"155 lbs.",
		"145 lbs.",
		"135 lbs.",
		"125 lbs.",
		"135 lbs",
		"125 lbs.",
		"115 lbs.",
	}

	response, err := http.Get(rankingsLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find(".wikitable").Each(func(rankingsTableIndex int, rankingsTable *goquery.Selection) {
		// Ignore legend table
		if rankingsTableIndex != 0 {

			var fighters []RankedFighter

			rankingsTable.Find("tr").Each(func(tableRowIndex int, tableRow *goquery.Selection) {
				if tableRowIndex >= 2 {
					fighter := RankedFighter{}
					fighter.Rank = strings.ReplaceAll(strings.TrimSpace(tableRow.Find("th").First().Text()), " (T)", "")

					tableRow.Find("td").Each(func(tableColumnIndex int, tableColumn *goquery.Selection) {
						if tableColumnIndex == 1 {
							fmt.Println("SCRAPING DATA FOR FIGHTER: " + strings.TrimSpace(tableColumn.Text()))

							fighter.Name = strings.TrimSpace(tableColumn.Text())

							fmt.Println("SCRAPING SHERDOG LINK FOR FIGHTER: " + strings.TrimSpace(tableColumn.Text()))

							fighter.SherdogLink = getSherdogLink(tableColumn.Find("a").First().AttrOr("href", ""))
						}

						var rankChangeIndex int

						if rankingsTableIndex <= 1 {
							rankChangeIndex = 4
						} else {
							rankChangeIndex = 3
						}

						if tableColumnIndex == rankChangeIndex {
							rankChangeSymbol := tableColumn.Find("img").First().AttrOr("alt", "Steady")

							if rankChangeSymbol == "New entry" {
								fighter.RankChange = "NR"
							} else if rankChangeSymbol == "Increase" || rankChangeSymbol == "Decrease" {
								fighter.RankChange = strings.TrimSpace(tableColumn.Text())
							} else if rankChangeSymbol == "Steady" {
								fighter.RankChange = "-"
							}
						}
					})
					fighters = append(fighters, fighter)
				}

			})

			ranking := Ranking{}

			ranking.Fighters = fighters
			ranking.WeightClass = weightClasses[rankingsTableIndex-1]
			ranking.Weight = weights[rankingsTableIndex-1]

			rankings = append(rankings, ranking)
		}
	})

	return rankings
}

func ScrapeUpcomingEvents() []Event {
	fmt.Println("SCRAPING UPCOMING EVENTS")

	events := []Event{}

	eventsLink := "https://en.wikipedia.org/wiki/List_of_UFC_events"

	response, err := http.Get(eventsLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	eventsTable := document.Find(".wikitable").First()

	var eventTableRows []*goquery.Selection

	eventsTable.Find("tr").Each(func(eventsTableRowIndex int, eventsTableRow *goquery.Selection) {
		eventTableRows = append(eventTableRows, eventsTableRow)
	})

	// Iterating in reverse for the sake of ordering, wikipedia puts the next fight at the bottom
	for i := len(eventTableRows) - 1; i >= 0; i-- {
		eventTableRows[i].Find("td").EachWithBreak(func(eventsTableColumnIndex int, eventsTableColumn *goquery.Selection) bool {
			if eventsTableColumnIndex == 0 {
				eventLink := eventsTableColumn.Find("a").First().AttrOr("href", "")
				if eventLink == "" {
					return false
				}

				if isFightScheduled(eventLink) == false {
					return false
				}

				events = append(events, ScrapeEventData(eventLink))
			}

			return true
		})
	}

	return events
}

func ScrapeEventData(eventLink string) Event {
	fmt.Println("SCRAPING EVENT")

	event := Event{}

	response, err := http.Get("https://en.wikipedia.org" + eventLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find(".toccolours").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(index int, tr *goquery.Selection) {
			if containsFightData(tr) {

			}
		})
	})

	return event
}

func ScrapeFighter() Fighter {
	fighter := Fighter{}

	return fighter
}

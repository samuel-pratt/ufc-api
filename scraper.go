package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Response struct {
	Rankings Rankings
	Events   []Event
}

type Rankings struct {
	MensPoundForPound   []RankedFighter
	WomensPoundForPound []RankedFighter
	Heavyweight         []RankedFighter
	LightHeavyweight    []RankedFighter
	Middleweight        []RankedFighter
	Welterweight        []RankedFighter
	Lightweight         []RankedFighter
	Featherweight       []RankedFighter
	Bantamweight        []RankedFighter
	Flyweight           []RankedFighter
	WomensBantamweight  []RankedFighter
	WomensFlyweight     []RankedFighter
	WomensStrawweight   []RankedFighter
}

type RankedFighter struct {
	Name       string
	Rank       string
	RankChange string
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

func ScrapeRankings() Rankings {
	rankings := Rankings{}

	rankingsLink := "https://en.wikipedia.org/wiki/UFC_Rankings"

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
					fighter.Rank = strings.TrimSpace(tableRow.Find("th").First().Text())

					tableRow.Find("td").Each(func(tableColumnIndex int, tableColumn *goquery.Selection) {
						if tableColumnIndex == 1 {
							fighter.Name = strings.TrimSpace(tableColumn.Text())
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

			// There's got to be a better way of doing this that I'm missing, but for now this will have to do
			// If you know a better way, open an issue suggesting it or a pr fixing it plz
			if rankingsTableIndex == 1 {
				rankings.MensPoundForPound = fighters
			} else if rankingsTableIndex == 2 {
				rankings.WomensPoundForPound = fighters
			} else if rankingsTableIndex == 3 {
				rankings.Heavyweight = fighters
			} else if rankingsTableIndex == 4 {
				rankings.LightHeavyweight = fighters
			} else if rankingsTableIndex == 5 {
				rankings.Middleweight = fighters
			} else if rankingsTableIndex == 6 {
				rankings.Welterweight = fighters
			} else if rankingsTableIndex == 7 {
				rankings.Lightweight = fighters
			} else if rankingsTableIndex == 8 {
				rankings.Featherweight = fighters
			} else if rankingsTableIndex == 9 {
				rankings.Bantamweight = fighters
			} else if rankingsTableIndex == 10 {
				rankings.Flyweight = fighters
			} else if rankingsTableIndex == 11 {
				rankings.WomensBantamweight = fighters
			} else if rankingsTableIndex == 12 {
				rankings.WomensFlyweight = fighters
			} else if rankingsTableIndex == 13 {
				rankings.WomensStrawweight = fighters
			}
		}
	})

	return rankings
}

func ScrapeUpcomingEvents() []Event {
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

	})

	return event
}

func ScrapeFighter() Fighter {
	fighter := Fighter{}

	return fighter
}

// Checks if the event at the given link has scheduled fights
func isFightScheduled(eventLink string) bool {
	result := false

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
		result = true
	})

	return result
}

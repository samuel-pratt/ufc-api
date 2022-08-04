package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Checks if a given row contains data from a fight
func containsFightData(tableRow *goquery.Selection) bool {
	result := true

	if strings.Contains(tableRow.Text(), "Main card") ||
		strings.Contains(tableRow.Text(), "Preliminary card") ||
		strings.Contains(tableRow.Text(), "Weight class") {
		result = false
	}

	return result
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

func getSherdogLink(fighterLink string) string {
	result := ""

	response, err := http.Get("https://en.wikipedia.org" + fighterLink)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	result = document.Find("a:contains('Professional MMA record for')").First().AttrOr("href", "")

	return result
}

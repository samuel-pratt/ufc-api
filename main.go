package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron"
)

var data Response

func updateData() {
	data = Scraper()

	fmt.Print("Updated data at: ")
	fmt.Println(time.Now())
}

func getData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonString, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func getRankings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonString, _ := json.Marshal(data.Rankings)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func getEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonString, _ := json.Marshal(data.Events)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func main() {
	// Create new schedule at startup
	updateData()

	// Schedule update every hour
	c := cron.New()
	c.AddFunc("@every 1h", updateData)
	c.Start()

	router := httprouter.New()

	// Routes
	router.GET("/api/", getData)
	router.GET("/api/rankings", getRankings)
	router.GET("/api/events", getEvents)

	var port = os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	http.ListenAndServe(":"+port, router)
}

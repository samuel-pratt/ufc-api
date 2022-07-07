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

func UpdateData() {
	data = Scraper()

	fmt.Print("Updated data at: ")
	fmt.Println(time.Now())
}

func GetOutages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonString, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func main() {
	// Create new schedule at startup
	UpdateData()

	// Schedule update every hour
	c := cron.New()
	c.AddFunc("@every 1h", UpdateData)
	c.Start()

	router := httprouter.New()

	// Root api call
	router.GET("/api/", GetOutages)

	var port = os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	http.ListenAndServe(":"+port, router)
}

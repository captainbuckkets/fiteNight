// Write these to be stored somewhere
// Print it out for now I guess

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getWebpage(pageURL string) {

	//var matchData []string

	// HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// Get our page - Method, Body, Error
	request, err := http.NewRequest("GET", pageURL, nil)
	// Error handling
	if err != nil {
		fmt.Println(err)
	}

	// Setting our headers
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("dnt", "1")
	request.Header.Set("upgrade-insecure-requests", "1")
	request.Header.Set("referrer", "https://www.betfair.com/*")
	// Send off our request
	resp, err := client.Do(request)

	// If we get a successful connection
	if resp.StatusCode == 200 {
		fmt.Println("Code 200")
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		// Error handling
		if err != nil {
			fmt.Println(err)
		}
		// Finds our row
		doc.Find("li .com-coupon-line-new-layout").Each(func(i int, d *goquery.Selection) {

			// Make an array for our match
			var match = make([]interface{}, 0)

			// Date | NameA | NameB | OddsA | OddsB | Snapshot Time

			// Find our date
			d.Find(".ui-countdown").Each(func(i int, dateEntry *goquery.Selection) {
				matchDate := dateEntry.Text()
				// Add our match to the array
				match = append(match, matchDate)
			})

			// Get our names
			d.Find(".team-name").Each(func(i int, nameEntry *goquery.Selection) {
				matchName := nameEntry.Text()
				// Add our name to the array
				match = append(match, matchName)
			})

			// Get our odds
			d.Find(".ui-display-decimal-price").Each(func(i int, oddEntry *goquery.Selection) {
				matchOdd := oddEntry.Text()
				// Add our odd to the array
				match = append(match, matchOdd)
			})

			// Get the current timestamp for graphs later
			now := time.Now()
			sec := now.Unix()
			match = append(match, sec)

			fmt.Println(match[0], match[1], match[2], match[3], match[4], match[5])

			// Write our new data to our MongoDB

			// Set client options
			clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

			// Connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)

		})
	}
}

func main() {
	getWebpage("https://www.betfair.com/sport/mixed-martial-arts")

}

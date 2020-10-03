// Write these to be stored somewhere
// Print it out for now I guess

package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define our structure for later use
type matchPoint struct {
	MatchDate string
	OpponentA string
	OpponentB string
	OddsA     float64
	OddsB     float64
	Timestamp int64
}

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

			// Date | NameA | NameB | OddsA | OddsB | Snapshot Time
			var date string
			var names []string
			var odds []float64
			var timeSnap int64

			// Find our date
			d.Find(".ui-countdown").Each(func(i int, dateEntry *goquery.Selection) {
				date = dateEntry.Text()
				date = strings.ReplaceAll(date, "\n", "")
			})

			// Get our names
			d.Find(".team-name").Each(func(i int, nameEntry *goquery.Selection) {
				matchName := nameEntry.Text()
				matchName = strings.ReplaceAll(matchName, "\n", "")
				// Add our name to the name array
				names = append(names, matchName)
			})

			var nameA = names[0]
			var nameB = names[1]

			// Get our odds
			d.Find(".ui-display-decimal-price").Each(func(i int, oddEntry *goquery.Selection) {
				odd := oddEntry.Text()
				// Remove the extra crap
				odd = strings.ReplaceAll(odd, "\n", "")
				oddFixed, err := strconv.ParseFloat(odd, 64)
				// Print our error if there is one
				if err != nil {
					fmt.Println(err)
				}
				// Add our odd to the array
				odds = append(odds, math.Round(oddFixed*100)/100)
			})

			var oddA = odds[0]
			var oddB = odds[1]
			// Get the current timestamp for graphs later
			now := time.Now()
			timeSnap = now.Unix()

			//fmt.Println(date, nameA, nameB, oddA, oddB, timeSnap)

			// Write our new data to our MongoDB
			clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
			// Connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)

			// Log our error if we have one
			if err != nil {
				log.Fatal(err)
			}

			// Check the connection to the database
			err = client.Ping(context.TODO(), nil)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Connected to MongoDB!")

			// Get our specific database and get our subccollection
			collection := client.Database("fiteNight").Collection("mma")

			newData := matchPoint{date, nameA, nameB, oddA, oddB, timeSnap}

			fmt.Println(newData)

			insertResult, err := collection.InsertOne(context.TODO(), newData)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Inserted a single document: ", insertResult.InsertedID)

			err = client.Disconnect(context.TODO())

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Connection to MongoDB closed.")

		})
	}
}

func main() {
	getWebpage("https://www.betfair.com/sport/mixed-martial-arts")

}

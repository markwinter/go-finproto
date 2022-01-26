/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/markwinter/go-finproto/itch"
)

func main() {
	filePath := flag.String("file", "", "Path to ITCH file")
	flag.Parse()

	config := itch.Configuration{
		MessageTypes: []byte{
			itch.MESSAGE_STOCK_DIRECTORY,
			itch.MESSAGE_PARTICIPANT_POSITION,
		},
		MaxMessages:    0,
		ReadBufferSize: itch.OneGB,
	}

	//defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	_, err := itch.ParseFile(*filePath, config)
	if err != nil {
		log.Fatal(err)
	}

	/*
		// If you have a lot of memory or a small file then read the whole file into memory
		data, err := os.ReadFile("01302020.NASDAQ_ITCH50")
		if err != nil {
			log.Fatal(err)
		}

		messages, err := itch.ParseMany(data, config)
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		// Print all the messages we passed
		for _, message := range messages {
			fmt.Printf("%v\n", message)
		}
	*/

	// Get participant positions for Goldman Sachs
	goldmanSachs := itch.MarketParticipants["GSCO"]
	fmt.Print(goldmanSachs[0]) // Print their first position

	// [Market Participant Position]
	// Stock Locate: 1176
	// Tracking Number: 0
	// Timestamp: 3h7m25.900657858s
	// MPID: GSCO
	// Stock: CAE
	// Primary: true
	// Mode: Normal
	// State: Active

	// Print the related stock for the position
	fmt.Print(itch.Directory[goldmanSachs[0].StockLocate])

	// [Stock Directory]
	// Stock Locate: 1176
	// Tracking Number: 0
	// Timestamp: 3h7m14.947382354s
	// Stock: CAE
	// Market Category: New York Stock Exchange (NYSE)
	// Financial Status Indicator: Not Available
	// Round Lot Size: 100
	// Round Lots Only: false
	// Issue Classification: Ordinary Share
	// Issue Sub-Type: Not Applicable
	// Authenticity: Live/Production
	// Short Sale Threshold Indicator: N
	// IPO Flag:
	// LULD Reference Price Tier: 2
	// ETP Flag: N
	// ETP Leverage Factor: 0
	// Inverse Indicator: false

	// Alternatively get stock by symbol using itch.StockMap
	stockLocate := itch.StockMap["AAPL"]
	fmt.Print(itch.Directory[stockLocate])

	// [Stock Directory]
	// Stock Locate: 13
	// Tracking Number: 0
	// Timestamp: 3h7m14.909934345s
	// Stock: AAPL
	// Market Category: Nasdaq Global Select Market
	// Financial Status Indicator: Normal
	// Round Lot Size: 100
	// Round Lots Only: false
	// Issue Classification: Common Stock
	// Issue Sub-Type: Not Applicable
	// Authenticity: Live/Production
	// Short Sale Threshold Indicator: N
	// IPO Flag: N
	// LULD Reference Price Tier: 1
	// ETP Flag: N
	// ETP Leverage Factor: 0
	// Inverse Indicator: false
}

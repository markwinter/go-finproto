package main

import (
	"log"
	"os"

	"github.com/markwinter/go-finproto/itch"
)

func main() {
	config := itch.Configuration{
		MessageTypes: []byte{},
		MaxMessages:  0,
	}

	/*
		messages, err := itch.ParseFile("01302020.NASDAQ_ITCH50", config)
		if err != nil {
			log.Fatal(err)
		}
	*/

	data, err := os.ReadFile("01302020.NASDAQ_ITCH50")
	if err != nil {
		log.Fatal(err)
	}

	// defer profile.Start().Stop()

	_, err = itch.ParseBytes(data, config)
	if err != nil {
		log.Fatal(err)
	}

	/*
		for _, message := range messages {
			fmt.Printf("%v\n", message)
		}
	*/
}

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
		MessageTypes:   []byte{'S'},
		MaxMessages:    0,
		ReadBufferSize: itch.OneGB,
	}

	//defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	messages, err := itch.ParseFile(*filePath, config)
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

	for _, message := range messages {
		fmt.Printf("%v\n", message)
	}
}

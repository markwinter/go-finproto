package main

import (
	"fmt"
	"log"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func main() {
	client := soupbintcp.Client{}
	client.Connect(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))

	err := client.Login("user", "pass")
	if err != nil {
		log.Printf("login failed: %v", err)
	}

	client.Logout()
	client.Disconnect()
}

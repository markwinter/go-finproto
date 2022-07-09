package main

import (
	"fmt"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func main() {
	client := soupbintcp.Client{}
	client.Connect(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))

	client.Login("user", "pass")
	client.Logout()

	client.Disconnect()
}

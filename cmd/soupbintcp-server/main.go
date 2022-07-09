package main

import (
	"fmt"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func main() {
	server := soupbintcp.Server{}
	server.ListenAndServe(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))
}

package main

import (
	"fmt"
	"log"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received packet: %v", packet)
}

func main() {
	client := soupbintcp.Client{
		PacketCallback: ReceivePacket,
	}
	client.Connect(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))

	err := client.Login("user", "pass")
	if err != nil {
		log.Printf("login failed: %v", err)
	}

	client.Receive()

	client.Logout()
	client.Disconnect()
}

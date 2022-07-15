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
		ServerAddr:     fmt.Sprintf("%s:%s", "127.0.0.1", "1337"),
		Username:       "user",
		Password:       "pass",
	}
	client.Connect()

	err := client.Login()
	if err != nil {
		log.Printf("login failed: %v", err)
	}

	client.SendDebugMessage("hello debug")
	client.Send([]byte("this is an unsequenced text message but could be bytes of a higher-level protocol packet"))

	client.Receive()

	client.Logout()
	client.Disconnect()
}

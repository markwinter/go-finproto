package main

import (
	"log"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received packet: %v", packet)
}

func main() {
	client := soupbintcp.Client{
		PacketCallback: ReceivePacket,
		ServerIp:       "127.0.0.1",
		ServerPort:     "1337",
		Username:       "test",
		Password:       "test",
	}

	if err := client.Login(); err != nil {
		log.Printf("login failed: %v\n", err)
		return
	}
	defer client.Logout()

	log.Println("logged in")

	if err := client.SendDebugMessage("hello debug"); err != nil {
		log.Printf("failed sending debug packet to server: %v\n", err)
	}

	if err := client.Send([]byte("this is an unsequenced text message but could be bytes of a higher-level protocol packet")); err != nil {
		log.Printf("failed sending packet to server: %v\n", err)
	}

	// Blocks until end of session packet received. Use a goroutine to unblock
	client.Receive()
}

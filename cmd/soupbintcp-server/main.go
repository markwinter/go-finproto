package main

import (
	"fmt"
	"log"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received unsequenced data packet: %v", packet)
}

func main() {
	server := soupbintcp.Server{
		// The Server can receive Unsequenced Data Packets from the Client
		PacketCallback: ReceivePacket,
	}
	server.ListenAndServe(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))
}

package main

import (
	"fmt"
	"log"

	"github.com/markwinter/go-finproto/soupbintcp"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received unsequenced data packet: %v", packet)
}

func LoginCallback(username, password string) bool {
	return username == "test" && password == "test"
}

func main() {
	server := soupbintcp.Server{
		// The Server can receive Unsequenced Data Packets from the Client
		PacketCallback: ReceivePacket,
		// All client login requests will invoke this callback
		LoginCallback: LoginCallback,
	}

	sessionId := "ABCDEFGHIJ"
	if err := server.CreateSession(sessionId); err != nil {
		log.Printf("failed to create session: %v", err)
		return
	}

	server.ListenAndServe(fmt.Sprintf("%s:%s", "127.0.0.1", "1337"))
}

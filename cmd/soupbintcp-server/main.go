package main

import (
	"log"

	soupbintcp "github.com/markwinter/go-finproto/soupbintcp/4.1"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received unsequenced data packet: %v\n", packet)
}

func DebugPacket(packet []byte) {
	log.Printf("Received debug packet: %v\n", packet)
}

func LoginCallback(username, password string) bool {
	return username == "test" && password == "test"
}

func main() {
	server := soupbintcp.NewServer(
		// All client login requests will invoke this callback
		soupbintcp.WithLoginCallback(LoginCallback),

		// The Server can receive Unsequenced Data Packets from the Client
		soupbintcp.WithPacketCallback(ReceivePacket),

		// Clients can send debug packets. Not normally used.
		soupbintcp.WithDebugCallback(DebugPacket),
	)

	sessionId := "ABCDEFGHIJ"
	if err := server.CreateSession(sessionId); err != nil {
		log.Printf("failed to create session: %v\n", err)
		return
	}

	server.ListenAndServe("127.0.0.1", "1337")
}

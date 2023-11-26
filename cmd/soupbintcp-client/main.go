package main

import (
	"log"

	soupbintcp "github.com/markwinter/go-finproto/soupbintcp/4.1"
)

func ReceivePacket(packet []byte) {
	log.Printf("Received packet: %v\n", packet)
}

func DebugPacket(text string) {
	log.Printf("[DEBUG] %s\n", text)
}

func main() {
	client := soupbintcp.NewClient(
		soupbintcp.WithServer("127.0.0.1", "4000"),
		soupbintcp.WithAuth("test", "test"),
		soupbintcp.WithCallback(ReceivePacket),
		soupbintcp.WithDebugCallback(DebugPacket),
	)

	// Login() will login with an empty session id and sequence number 1, as is recommended by the spec
	// If you know the session and sequence number, use LoginSession()
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
	// Will automatically attempt to re-connect and resume the same session and sequence number
	client.Receive()
}

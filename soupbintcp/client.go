package soupbintcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	heartbeatPeriod = 1000
)

type Client struct {
	PacketCallback    func([]byte)
	conn              net.Conn
	sequenceNumber    int
	session           string
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan bool
}

func (c *Client) Connect(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	c.conn = conn
}

func (c *Client) Disconnect() {
	c.conn.Close()
}

func (c *Client) Login(username, password string) error {
	username = fmt.Sprintf("%-6s", username)
	password = fmt.Sprintf("%-10s", password)

	request := LoginRequestPacket{
		Packet: Packet{
			Length: [2]byte{0, 52},
			Type:   'L',
		},
	}

	copy(request.Username[:], username)
	copy(request.Password[:], password)
	copy(request.HeartbeatTimeout[:], fmt.Sprint(heartbeatPeriod))
	copy(request.SequenceNumber[:], "1")
	copy(request.Session[:], "          ")

	binary.Write(c.conn, binary.BigEndian, &request)

	packet, err := GetNextPacket(c.conn)
	if err != nil {
		return err
	}

	switch packet[2] {
	case 'A':
		c.handleLoginAccepted(packet)
	case 'J':
		switch packet[3] {
		case 'A':
			return errors.New("not authorized")
		case 'S':
			return errors.New("session not available")
		}
	}

	go c.runHeartbeat()

	return nil
}

func (c *Client) Logout() {
	request := LogoutRequestPacket{
		Packet: Packet{
			Length: [2]byte{0, 1},
			Type:   'O',
		},
	}

	binary.Write(c.conn, binary.BigEndian, &request)

	c.heartbeatStopChan <- true
}

func (c *Client) runHeartbeat() {
	c.heartbeatTicker = time.NewTicker(heartbeatPeriod * time.Millisecond)
	c.heartbeatStopChan = make(chan bool)

	for {
		select {
		case <-c.heartbeatTicker.C:
			c.sendHeartbeat()
		case <-c.heartbeatStopChan:
			return
		}
	}
}

func (c *Client) sendHeartbeat() {
	request := HeartbeatPacket{
		Packet: Packet{
			Length: [2]byte{0, 3},
			Type:   'R',
		},
	}

	binary.Write(c.conn, binary.BigEndian, &request)
}

func (c *Client) Receive() {
	for {
		// We give a grace period of 2 * heartbeatPeriod to read something from the server
		// The server should be sending heartbeats every 1 * heartbeatPeriod
		c.conn.SetReadDeadline(time.Now().Add(heartbeatPeriod * time.Millisecond * 2))

		packet, err := GetNextPacket(c.conn)
		if err != nil {
			return
		}

		switch packet[2] {
		case '+':
			handleDebugPacket(packet)
			sendDebugPacket("pong", c.conn)
		case 'H':
			log.Print("received heartbeat packet")
		case 'Z':
			log.Print("end of session packet")
			return
		case 'S':
			c.sequenceNumber++
			c.PacketCallback(packet[3:])
		case 'U':
			c.PacketCallback(packet[3:])
		default:
			log.Print("unknown packet type received")
		}
	}
}

func (c *Client) handleLoginAccepted(packet []byte) {
	c.session = strings.TrimSpace(string(packet[3:13]))

	sq := strings.TrimSpace(string(packet[13:33]))
	seq, err := strconv.Atoi(sq)
	if err != nil {
		log.Print(err)
	}
	c.sequenceNumber = seq
}

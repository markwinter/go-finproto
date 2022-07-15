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
	ServerAddr        string
	Username          string
	Password          string
	conn              net.Conn
	sequenceNumber    int
	session           string
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan bool
}

func (c *Client) Connect() {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		log.Panic(err)
	}
	c.conn = conn
}

func (c *Client) Disconnect() {
	c.conn.Close()
}

// Login will try to connect to a new session and start receiving the first message
// If you want to connect to a specific session then use LoginSession
func (c *Client) Login() error {
	username := fmt.Sprintf("%-6s", c.Username)
	password := fmt.Sprintf("%-10s", c.Password)

	request := LoginRequestPacket{
		Packet: Packet{
			Length: [2]byte{0, 52},
			Type:   'L',
		},
	}

	copy(request.Username[:], username)
	copy(request.Password[:], password)
	copy(request.HeartbeatTimeout[:], fmt.Sprint(heartbeatPeriod))
	copy(request.SequenceNumber[:], fmt.Sprintf("%20s", "1"))
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

// LoginSession is used to connect to a known existing session
// and start receiving from the given sequence number
func (c *Client) LoginSession(session, sequence string) error {
	username := fmt.Sprintf("%-6s", c.Username)
	password := fmt.Sprintf("%-10s", c.Password)

	request := LoginRequestPacket{
		Packet: Packet{
			Length: [2]byte{0, 52},
			Type:   'L',
		},
	}

	copy(request.Username[:], username)
	copy(request.Password[:], password)
	copy(request.HeartbeatTimeout[:], fmt.Sprint(heartbeatPeriod))
	copy(request.SequenceNumber[:], fmt.Sprintf("%20s", sequence))
	copy(request.Session[:], fmt.Sprintf("%10s", session))

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

// Receive will start listening for packets from the Server. Any sequenced or unsequenced data
// packets will be passed to your PacketCallback function. Receive will also attempt to automatically
// reconnect to the Server and rejoin the previous session with the current message sequence number.
// Receive will block until an end of session packet is received.
func (c *Client) Receive() {
	for {
		// We give a grace period of 2 * heartbeatPeriod to read something from the server
		// The server should be sending heartbeats every 1 * heartbeatPeriod
		c.conn.SetReadDeadline(time.Now().Add(heartbeatPeriod * time.Millisecond * 2))

		packet, err := GetNextPacket(c.conn)
		if err != nil {
			// Try to reconnect and rejoin previous session with current sequenceNumber
			c.Connect()
			c.LoginSession(c.session, fmt.Sprint(c.sequenceNumber))
			continue
		}

		switch packet[2] {
		case '+':
			handleDebugPacket(packet)
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

// Send an unsequenced data packet to the Server
func (c *Client) Send(data []byte) {
	l := uint16(2 + 1 + len(data))
	buf := make([]byte, l)

	binary.BigEndian.PutUint16(buf[0:2], l-2)
	buf[2] = 'U'
	copy(buf[3:], data)

	binary.Write(c.conn, binary.BigEndian, &buf)
}

// Send a debug packet with human readable text to the Server. Not normally used.
func (c *Client) SendDebugMessage(text string) {
	sendDebugPacket(text, c.conn)
}

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

type Client struct {
	// PacketCallback is called for every sequenced packet received. The byte slice parameter contains just the message data and should be further
	// parsed as some higher level protocol
	PacketCallback func([]byte)
	// UnsequencedCallback is called for every unsequenced packet received. The byte slice parameter contains just the message data and should be further
	// parsed as some higher level protocol
	UnsequencedCallback func([]byte)
	ServerIp            string
	ServerPort          string
	Username            string
	Password            string

	conn              net.Conn
	sequenceNumber    uint64
	session           string
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan bool
	sentMessageChan   chan bool
}

type ClientOption func(client *Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.Username = username
		c.Password = password
	}
}

func WithCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.PacketCallback = callback
	}
}

func WithUnsequencedCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.UnsequencedCallback = callback
	}
}

func WithServer(ip, port string) ClientOption {
	return func(c *Client) {
		c.ServerIp = ip
		c.ServerPort = port
	}
}

func (c *Client) connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.ServerIp, c.ServerPort))
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) disconnect() {
	c.conn.Close()
}

// Login will try to connect to a new session and start receiving the first message
// If you want to connect to a specific session then use LoginSession
func (c *Client) Login() error {
	return c.login("", 1)
}

// LoginSession is used to connect to a known existing session
// and start receiving from the given sequence number
func (c *Client) LoginSession(session string, sequence uint64) error {
	return c.login(session, sequence)
}

func (c *Client) login(session string, sequence uint64) error {
	if err := c.connect(); err != nil {
		return err
	}

	loginPacket := makeLoginRequestPacket(c.Username, c.Password, session, sequence)

	if err := binary.Write(c.conn, binary.BigEndian, loginPacket.Bytes()); err != nil {
		return err
	}

	packet, err := getNextPacket(c.conn)
	if err != nil {
		return err
	}

	switch packet[2] {
	case PacketTypeLoginAccepted:
		if err := c.handleLoginAccepted(packet); err != nil {
			return err
		}
	case PacketTypeLoginRejected:
		switch packet[3] {
		case LoginRejectedNotAuthorized:
			return errors.New("not authorized")
		case LoginRejectedSessionUnavailable:
			return errors.New("session not available")
		}
	}

	go c.runHeartbeat()

	return nil
}

// Logout from the Server
func (c *Client) Logout() {
	packet := makeLogoutRequestPacket()

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		log.Println("failed sending logout request")
	}

	c.heartbeatStopChan <- true

	c.disconnect()
}

func (c *Client) runHeartbeat() {
	c.sentMessageChan = make(chan bool)
	c.heartbeatStopChan = make(chan bool)

	c.heartbeatTicker = time.NewTicker(heartbeatPeriod_ms * time.Millisecond)
	defer c.heartbeatTicker.Stop()

	for {
		select {
		case <-c.heartbeatTicker.C:
			c.sendHeartbeat()
		case <-c.sentMessageChan:
			// If we sent a message to the server, reset ticker so we're not sending
			// more often than we need to
			c.heartbeatTicker.Reset(heartbeatPeriod_ms * time.Millisecond)
			continue
		case <-c.heartbeatStopChan:
			return
		}
	}
}

func (c *Client) sendHeartbeat() {
	packet := makeClientHeartbeatPacket()

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		log.Println("failed sending heartbeat")
	}
}

// Receive will start listening for packets from the Server. Any sequenced or unsequenced data
// packets will be passed to your PacketCallback function. Receive will also attempt to automatically
// reconnect to the Server and rejoin the previous session with the current message sequence number.
// Receive will block until an end of session packet is received.
func (c *Client) Receive() {
	for {
		// We give a grace period of 2 * heartbeatPeriod to read something from the server
		// The server should be sending heartbeats every 1 * heartbeatPeriod
		err := c.conn.SetReadDeadline(time.Now().Add(heartbeatPeriod_ms * time.Millisecond * 2))
		if err != nil {
			log.Println("error setting read deadline")
			continue
		}

		packet, err := getNextPacket(c.conn)
		if err != nil {
			log.Printf("connection error, attempting to relogin to session %q with sequence number %d\n", c.session, c.sequenceNumber)
			// Try to reconnect and rejoin previous session with current sequenceNumber
			// TODO: some exponential backoff and max retry logic
			if err := c.LoginSession(c.session, c.sequenceNumber); err != nil {
				log.Println("failed to login after reconnect")
				return
			}

			continue
		}

		switch packet[2] {
		case PacketTypeDebug:
			log.Printf("[DEBUG PACKET] %s", packet[3:])
		case PacketTypeServerHeartbeat:
			log.Print("received heartbeat packet")
		case PacketTypeEndOfSession:
			log.Print("end of session packet")
			return
		case PacketTypeSequencedData:
			c.sequenceNumber++
			if c.PacketCallback != nil {
				c.PacketCallback(packet[3:])
			}
		case PacketTypeUnsequencedData:
			if c.UnsequencedCallback != nil {
				c.UnsequencedCallback(packet[3:])
			}
		default:
			log.Print("unknown packet type received")
		}
	}
}

func (c *Client) handleLoginAccepted(packet []byte) error {
	c.session = strings.TrimSpace(string(packet[3:13]))

	sq := strings.TrimSpace(string(packet[13:33]))

	var err error
	c.sequenceNumber, err = strconv.ParseUint(sq, 10, 64)
	if err != nil {
		return err
	}

	log.Printf("connected to session %q and starting with sequence %d\n", c.session, c.sequenceNumber)

	return nil
}

// Send an unsequenced data packet to the Server
func (c *Client) Send(data []byte) error {
	packet := makeUnsequencedDataPacket(data)

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		return err
	}

	c.sentMessageChan <- true
	return nil
}

// Send a debug packet with human readable text to the Server. Not normally used.
func (c *Client) SendDebugMessage(text string) error {
	packet := makeDebugPacket(text)

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		return err
	}

	c.sentMessageChan <- true
	return nil
}

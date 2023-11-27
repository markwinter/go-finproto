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
	// packetCallback is called for every sequenced packet received. The byte slice parameter contains just the message data and should be further
	// parsed as some higher level protocol
	packetCallback func([]byte)
	// unsequencedCallback is called for every unsequenced packet received. The byte slice parameter contains just the message data and should be further
	// parsed as some higher level protocol
	unsequencedCallback func([]byte)
	// debugCallback is called for every debug packet received. This is not normally used
	debugCallBack func(string)

	serverIp   string
	serverPort string
	username   string
	password   string

	conn              net.Conn
	sequenceNumber    uint64
	session           string
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan bool
	sentMessageChan   chan bool
}

type ClientOption func(client *Client)

// NewClient creates a new soupbintcp client. The default parameters can be configured using ClientOptions passed in as parameters
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		session:        "",
		sequenceNumber: 0, // 0 indicates start receiving most recently generated message

		sentMessageChan:   make(chan bool),
		heartbeatStopChan: make(chan bool),
		heartbeatTicker:   time.NewTicker(heartbeatPeriod_ms * time.Millisecond),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithAuth sets the username and password to use when connecting to the Server
func WithAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// WithCallback sets the callback function for every sequenced packet received
func WithCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.packetCallback = callback
	}
}

// WithUnsequencedCallback sets the callback function for every unsequenced packet received
func WithUnsequencedCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.unsequencedCallback = callback
	}
}

// WithDebugCallback sets the callback function for every debug packet received. Not normally used
func WithDebugCallback(callback func(string)) ClientOption {
	return func(c *Client) {
		c.debugCallBack = callback
	}
}

// WithServer sets the ip and port to use to connect to the Server
func WithServer(ip, port string) ClientOption {
	return func(c *Client) {
		c.serverIp = ip
		c.serverPort = port
	}
}

// WithSession sets the initial session and sequence number when connecting to the Server
func WithSession(id string, sequence uint64) ClientOption {
	return func(c *Client) {
		c.session = id
		c.sequenceNumber = sequence
	}
}

func (c *Client) connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.serverIp, c.serverPort))
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
// If you want to connect to a specific session, use WithSession() in the NewClient options
func (c *Client) Login() error {
	if err := c.connect(); err != nil {
		return err
	}

	loginPacket := makeLoginRequestPacket(c.username, c.password, c.session, c.sequenceNumber)

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

// Logout will logout from the Server
func (c *Client) Logout() {
	packet := makeLogoutRequestPacket()

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		log.Printf("failed sending logout request: %v\n", err)
	}

	c.heartbeatStopChan <- true

	c.disconnect()
}

func (c *Client) runHeartbeat() {
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
		log.Printf("failed sending heartbeat: %v\n", err)
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
			if err := c.Login(); err != nil {
				log.Println("failed to login after reconnect")
				return
			}

			continue
		}

		switch packet[2] {
		case PacketTypeDebug:
			if c.debugCallBack != nil {
				c.debugCallBack(string(packet[3:]))
			}
		case PacketTypeServerHeartbeat:
			log.Print("received heartbeat packet")
		case PacketTypeEndOfSession:
			log.Print("end of session packet")
			return
		case PacketTypeSequencedData:
			c.sequenceNumber++
			if c.packetCallback != nil {
				c.packetCallback(packet[3:])
			}
		case PacketTypeUnsequencedData:
			if c.unsequencedCallback != nil {
				c.unsequencedCallback(packet[3:])
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

// Send sends an unsequenced data packet to the Server
func (c *Client) Send(data []byte) error {
	packet := makeUnsequencedDataPacket(data)

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		return err
	}

	c.sentMessageChan <- true
	return nil
}

// SendDebugMessage sends a debug packet with human readable text to the Server. Not normally used.
func (c *Client) SendDebugMessage(text string) error {
	packet := makeDebugPacket(text)

	if err := binary.Write(c.conn, binary.BigEndian, packet.Bytes()); err != nil {
		return err
	}

	c.sentMessageChan <- true
	return nil
}

// CurrentSession returns the current session id
func (c *Client) CurrentSession() string {
	return c.session
}

// CurrentSequenceNumber returns the client's current sequence number
func (c *Client) CurrentSequenceNumber() uint64 {
	return c.sequenceNumber
}

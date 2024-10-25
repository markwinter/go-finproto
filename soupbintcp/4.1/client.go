package soupbintcp

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
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

	serverAddr string

	username string
	password string

	conn               net.Conn
	compressionEnabled bool

	sequenceNumber uint64
	session        string

	heartbeatStopChan chan bool
	sentMessageChan   chan bool

	backoff *backoff.ExponentialBackOff
}

// NewClient creates a new soupbintcp client. The default parameters can be configured using ClientOptions passed in as parameters
func NewClient(addr string, opts ...ClientOption) *Client {
	b := backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(100*time.Millisecond),
		backoff.WithMaxElapsedTime(30*time.Second),
		backoff.WithMaxInterval(5*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithRandomizationFactor(0.1),
	)

	c := &Client{
		serverAddr: addr,

		session:        "",
		sequenceNumber: 0, // 0 indicates start receiving most recently generated message

		sentMessageChan:   make(chan bool),
		heartbeatStopChan: make(chan bool),

		backoff: b,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) connect() error {
	conn, err := net.Dial("tcp", c.serverAddr)
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
	ticker := time.NewTicker(heartbeatPeriod_ms * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.sendHeartbeat()
		case <-c.sentMessageChan:
			// If we sent a message to the server, reset ticker so we're not sending
			// more often than we need to
			ticker.Reset(heartbeatPeriod_ms * time.Millisecond)
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

func (c *Client) reconnect() error {
	c.heartbeatStopChan <- true

	notify := func(err error, t time.Duration) {
		log.Printf("retrying connetion in: %s\n", t)
	}

	log.Printf("connection error, attempting to relogin to session %q with sequence number %d\n", c.session, c.sequenceNumber)

	if err := backoff.RetryNotify(c.Login, c.backoff, notify); err != nil {
		log.Println("failed to reconnect to the server after max retries")
		return err
	}

	return nil
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
			if err := c.reconnect(); err != nil {
				return
			}
			continue
		}

		packet, err := getNextPacket(c.conn)
		if err != nil {
			log.Printf("error getting packet: %v", err)
			if err := c.reconnect(); err != nil {
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

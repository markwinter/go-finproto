package soupbintcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn           net.Conn
	sequenceNumber int
	session        string
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
	copy(request.HeartbeatTimeout[:], "1000")
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

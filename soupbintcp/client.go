package soupbintcp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
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

func (c *Client) Login(username, password string) {
	username = fmt.Sprintf("%6s", username)
	password = fmt.Sprintf("%10s", password)

	request := LoginRequestPacket{
		Type: 'L',
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, 52)
	copy(request.Length[:], buf)

	copy(request.Username[:], username)
	copy(request.Password[:], password)
	copy(request.HeartbeatTimeout[:], "1000")
	copy(request.SequenceNumber[:], "1")

	log.Print(request)

	binary.Write(c.conn, binary.BigEndian, &request)
}

func (c *Client) Logout() {
	request := LogoutRequestPacket{
		Type: 'O',
	}

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, 1)
	copy(request.Length[:], buf)

	log.Print(request)

	binary.Write(c.conn, binary.BigEndian, &request)
}

type LoginRequestPacket struct {
	Length           [2]byte
	Type             byte
	Username         [6]byte
	Password         [10]byte
	Session          [10]byte
	SequenceNumber   [20]byte
	HeartbeatTimeout [5]byte
}

type LogoutRequestPacket struct {
	Length [2]byte
	Type   byte
}

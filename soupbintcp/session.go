package soupbintcp

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

type Session struct {
	io.WriteCloser

	id             string
	sequenceNumber uint64

	conn net.Conn

	stopChan chan bool
}

func MakeSession(id string, conn net.Conn) Session {
	session := Session{
		id:             id,
		sequenceNumber: 1, // The sequence number of the first sequenced message in each session is always 1
		conn:           conn,
		stopChan:       make(chan bool),
	}

	go session.run()

	return session
}

func (s *Session) Write(p []byte) (n int, err error) {
	if err := binary.Write(s.conn, binary.BigEndian, p); err != nil {
		return 0, err
	}
	s.sequenceNumber += 1
	return len(p), nil
}

func (s *Session) Close() error {
	s.stopChan <- true
	return s.conn.Close()
}

func (s *Session) run() {
	heartbeatTicker := time.NewTicker(heartbeatPeriod_ms * time.Millisecond)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			s.sendHeartbeat()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Session) sendHeartbeat() {
	request := HeartbeatPacket{
		Packet: Packet{
			Length: [2]byte{0, 3},
			Type:   PacketServerHeartbeat,
		},
	}

	if err := binary.Write(s.conn, binary.BigEndian, &request); err != nil {
		log.Printf("failed sending heartbeat packet: %v\n", err)
	}
}

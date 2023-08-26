package soupbintcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	LoginCallback  func(LoginRequestPacket) bool
	PacketCallback func([]byte)

	activeSession bool
	session       Session
}

func (s *Server) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	log.Printf("Listening on %s", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		log.Printf("client connected %q", conn.RemoteAddr())

		// Get Login Request Packet
		// Check session is active
		// Check sequence number
		// Start sending from that sequence number

		s.sendLoginAccepted(conn)
		s.session.AddConn(conn)
	}
}

func (s *Server) CreateSession(id string) error {
	if s.activeSession {
		return errors.New("session already exists, call DeleteSession first")
	}

	s.session = MakeSession(id)
	s.activeSession = true

	return nil
}

func (s *Server) DeleteSession(id string) error {
	if !s.activeSession {
		return nil
	}

	// TODO: call end of session packet here?
	return nil
}

func (s *Server) SendToSession(data []byte) error {
	if !s.activeSession {
		return errors.New("no active session")
	}

	return s.session.Send(data)
}

func (s *Server) sendLoginAccepted(conn net.Conn) {
	request := LoginAcceptedPacket{
		Packet: Packet{
			Length: [2]byte{0, 33},
			Type:   PacketLoginAccepted,
		},
	}

	copy(request.SequenceNumber[:], fmt.Sprintf("%20d", s.session.CurrentSequenceNumber))
	copy(request.Session[:], []byte(s.session.Id))

	if err := binary.Write(conn, binary.BigEndian, &request); err != nil {
		log.Printf("failed sending login accepted: %v", err)
	}
}

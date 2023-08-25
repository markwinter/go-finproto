package soupbintcp

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionIdLength = 10
)

type Server struct {
	LoginCallback  func(LoginRequestPacket) bool
	PacketCallback func([]byte)
	sequenceNumber uint64
	sessions       map[string]Session
	sentData       [][]byte
}

func (s *Server) ListenAndServe(addr string) {
	s.sessions = map[string]Session{}

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

		// Ensure we generate an id we haven't already used
		// Todo: make this smarter, possible infinite loop if we're unlucky
		id := generateSessionId(sessionIdLength)

		s.sendLoginAccepted(id, conn)

		for _, ok := s.sessions[id]; ok; _, ok = s.sessions[id] {
			id = generateSessionId(sessionIdLength)
		}

		session := MakeSession(id, conn)
		s.sessions[id] = session
	}
}

func (s *Server) SendData(data []byte) error {
	s.sentData = append(s.sentData, data)
	s.sequenceNumber += 1
	return nil
}

func generateSessionId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *Server) sendLoginAccepted(session string, conn net.Conn) {
	request := LoginAcceptedPacket{
		Packet: Packet{
			Length: [2]byte{0, 33},
			Type:   'A',
		},
	}

	copy(request.SequenceNumber[:], fmt.Sprintf("%20d", s.sequenceNumber))
	copy(request.Session[:], []byte(session))

	log.Print(request)

	if err := binary.Write(conn, binary.BigEndian, &request); err != nil {
		log.Printf("failed sending login accepted: %v", err)
	}
}

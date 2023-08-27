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

type Server struct {
	LoginCallback  func(username, password string) bool
	PacketCallback func([]byte)

	activeSession bool
	session       session
}

func (s *Server) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	log.Printf("Listening on %s\n", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		log.Printf("client connected %q\n", conn.RemoteAddr())

		if !s.activeSession {
			s.sendLoginRejected(LoginRejectedSessionUnavailable, conn)
			return
		}

		go s.handleConnection(conn)
	}
}

// CreateSession creates a new session. There can only be one active session at a time
func (s *Server) CreateSession(id string) error {
	if s.activeSession {
		return errors.New("session already exists, call DeleteSession first")
	}

	s.session = makeSession(id)
	s.activeSession = true

	return nil
}

// DeleteSession deletes an active session
func (s *Server) DeleteSession(id string) error {
	if !s.activeSession {
		return nil
	}

	// TODO: call end of session packet here?
	s.activeSession = false
	return nil
}

// SendToSession adds data to the session that will eventually be broadcast to all clients connected to the session.
// A session must first have been created with CreateSession
func (s *Server) SendToSession(data []byte) error {
	if !s.activeSession {
		return errors.New("no active session")
	}

	if err := s.session.dataStore.store(data); err != nil {
		return err
	}
	s.session.currentSequenceNumber++
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	if !s.handleLogin(conn) {
		return
	}

	// Start sending from sequence number
	// Listen for unsequenced and debug packets
}

func (s *Server) handleLogin(conn net.Conn) bool {
	// Clients should immediately be sending Login Requests after establishing the connection
	// so set a deadline of 2 heartbeats to receive the request
	if err := conn.SetReadDeadline(time.Now().Add(heartbeatPeriod_ms * time.Millisecond * 2)); err != nil {
		log.Println("error setting read deadline")
		return false
	}

	packet, err := getNextPacket(conn)
	if err != nil || packet[2] != PacketLoginRequest {
		s.sendLoginRejected(LoginRejectedNotAuthorized, conn)
		return false
	}

	username := strings.TrimSpace(string(packet[3:9]))
	password := strings.TrimSpace(string(packet[9:19]))

	if s.LoginCallback != nil && !s.LoginCallback(username, password) {
		s.sendLoginRejected(LoginRejectedNotAuthorized, conn)
		return false
	}

	requestedSession := strings.TrimSpace(string(packet[19:29]))
	if requestedSession != "" && requestedSession != s.session.id {
		s.sendLoginRejected(LoginRejectedSessionUnavailable, conn)
		return false
	}

	seq := strings.TrimSpace(string(packet[29:49]))
	requestedSeq, err := strconv.ParseUint(seq, 10, 64)
	if err != nil {
		log.Printf("failed to parse requested sequence number: %v\n", err)
		return false
	}

	startSeq := requestedSeq
	if requestedSeq == 0 || requestedSeq > s.session.currentSequenceNumber {
		startSeq = s.session.currentSequenceNumber
	}

	log.Printf("starting seq for client %q is %d", conn.RemoteAddr(), startSeq)

	s.sendLoginAccepted(startSeq, conn)

	return true
}

func (s *Server) sendLoginAccepted(seq uint64, conn net.Conn) {
	request := LoginAcceptedPacket{
		Header: Header{
			Length: [2]byte{0, 33},
			Type:   PacketLoginAccepted,
		},
	}

	copy(request.SequenceNumber[:], fmt.Sprintf("%20d", seq))
	copy(request.Session[:], []byte(s.session.id))

	if err := binary.Write(conn, binary.BigEndian, &request); err != nil {
		log.Printf("failed sending login accepted: %v\n", err)
	}
}

func (s *Server) sendLoginRejected(reason byte, conn net.Conn) {
	request := LoginRejectedPacket{
		Header: Header{
			Length: [2]byte{0, 33},
			Type:   PacketLoginRejected,
		},
		Reason: reason,
	}

	if err := binary.Write(conn, binary.BigEndian, &request); err != nil {
		log.Printf("failed sending login rejected: %v\n", err)
	}
}

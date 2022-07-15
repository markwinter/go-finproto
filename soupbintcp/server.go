package soupbintcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	PacketCallback    func([]byte)
	activeConns       sync.Map
	sequenceNumber    uint64
	session           string
	heartbeatTicker   *time.Ticker
	heartbeatStopChan chan bool
}

func (s *Server) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	log.Printf("Listening on %s", addr)

	s.sequenceNumber = 1
	s.session = "abcdefghij"

	go s.runHeartbeat()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) closeConn(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	s.activeConns.Delete(remoteAddr)
	conn.Close()
}

func (s *Server) sendHeartbeat() {
	request := HeartbeatPacket{
		Packet: Packet{
			Length: [2]byte{0, 3},
			Type:   'H',
		},
	}

	s.activeConns.Range(func(key, conn any) bool {
		binary.Write(conn.(net.Conn), binary.BigEndian, &request)
		return true
	})
}

func (s *Server) runHeartbeat() {
	s.heartbeatTicker = time.NewTicker(heartbeatPeriod * time.Millisecond)
	s.heartbeatStopChan = make(chan bool)

	for {
		select {
		case <-s.heartbeatTicker.C:
			s.sendHeartbeat()
		case <-s.heartbeatStopChan:
			return
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	s.activeConns.Store(remoteAddr, conn)

	for {
		// We give a grace period of 2 * heartbeatPeriod to read something from the client
		// The client should be sending heartbeats every 1 * heartbeatPeriod
		conn.SetReadDeadline(time.Now().Add(heartbeatPeriod * time.Millisecond * 2))

		packet, err := GetNextPacket(conn)
		if err != nil {
			s.closeConn(conn)
			log.Print(err)
			return
		}

		switch packet[2] {
		case '+':
			handleDebugPacket(packet)
		case 'L':
			log.Printf("Login received from %s", conn.RemoteAddr().String())
			s.handleLogin(conn, packet)
		case 'O':
			log.Printf("Logout received from %s", conn.RemoteAddr().String())
			s.handleLogout(conn)
			return
		case 'R':
			log.Printf("Received heartbeat from: %s", conn.RemoteAddr().String())
		case 'U':
			s.PacketCallback(packet[3:])
		}
	}
}

func (s *Server) handleLogin(conn net.Conn, data []byte) {
	username := strings.TrimSpace(string(data[3:9]))
	password := strings.TrimSpace(string(data[9:19]))
	log.Printf("username: %s password: %s", username, password)

	session := strings.TrimSpace(string(data[19:29]))
	log.Printf("session: %v", session)

	if session != "" && session != s.session {
		s.sendLoginRejected(conn, 'S')
	}

	sequence := strings.TrimSpace(string(data[29:49]))
	log.Printf("sequence: %v", sequence)

	ht := string(bytes.Trim(data[49:], "\x00"))
	t, err := strconv.Atoi(ht)
	if err != nil {
		log.Print(err)
		t = 2000
	}
	heartbeat := time.Millisecond * time.Duration(t)
	log.Printf("heartbeat: %v", heartbeat.Milliseconds())

	// We will just accept all logins.
	// Later may want to add proper auth.
	s.sendLoginAccepted(conn)
}

func (s *Server) handleLogout(conn net.Conn) {
	s.closeConn(conn)
}

func (s *Server) sendLoginAccepted(conn net.Conn) {
	request := LoginAcceptedPacket{
		Packet: Packet{
			Length: [2]byte{0, 33},
			Type:   'A',
		},
	}

	copy(request.SequenceNumber[:], fmt.Sprintf("%20d", s.sequenceNumber))
	copy(request.Session[:], []byte(s.session))

	log.Print(request)

	binary.Write(conn, binary.BigEndian, &request)
}

func (s *Server) sendLoginRejected(conn net.Conn, reason byte) {
	request := LoginRejectedPacket{
		Packet: Packet{
			Length: [2]byte{0, 4},
			Type:   'J',
		},
		Reason: reason,
	}

	binary.Write(conn, binary.BigEndian, &request)
}

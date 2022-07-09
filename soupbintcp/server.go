package soupbintcp

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct{}

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

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	packetLengthBuffer := make([]byte, 2)
	for {
		_, err := conn.Read(packetLengthBuffer)
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			return
		}
		packetLength := binary.BigEndian.Uint16(packetLengthBuffer)
		log.Printf("length: %d", packetLength)

		buf := make([]byte, packetLength+2)
		copy(buf[0:2], packetLengthBuffer)
		_, err = conn.Read(buf[2:])
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			return
		}

		switch buf[2] {
		case 'L':
			log.Printf("Login received from %s", conn.RemoteAddr().String())
			handleLogin(buf)
		case 'O':
			log.Printf("Logout received from %s", conn.RemoteAddr().String())
			conn.Close()
			return
		}
	}
}

func handleLogin(data []byte) {
	log.Printf("%d %v", len(data), data)

	username := strings.TrimSpace(string(data[3:9]))
	password := strings.TrimSpace(string(data[9:19]))

	log.Printf("username: %s password: %s", username, password)

	session := strings.TrimSpace(string(data[19:29]))
	log.Printf("session: %v", session)

	sequence := strings.TrimSpace(string(data[29:49]))
	log.Printf("sequence: %v", sequence)

	ht := strings.TrimSpace(string(data[49:53]))
	t, err := strconv.Atoi(ht)
	if err != nil {
		log.Print(err)
		t = 2000
	}
	heartbeat := time.Millisecond * time.Duration(t)
	log.Printf("heartbeat: %v", heartbeat.Milliseconds())
}

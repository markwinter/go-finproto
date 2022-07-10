package soupbintcp

import (
	"encoding/binary"
	"log"
	"net"
)

func GetNextPacket(conn net.Conn) ([]byte, error) {
	packetLengthBuffer := make([]byte, 2)
	_, err := conn.Read(packetLengthBuffer)
	if err != nil {
		log.Printf("Error reading: %v\n", err)
		return []byte{}, err
	}
	packetLength := binary.BigEndian.Uint16(packetLengthBuffer)

	buf := make([]byte, packetLength+2)
	copy(buf[0:2], packetLengthBuffer)
	_, err = conn.Read(buf[2:])
	if err != nil {
		log.Printf("Error reading: %v\n", err)
		return []byte{}, err
	}

	return buf, nil
}

func sendDebugPacket(text string, conn net.Conn) {
	l := uint16(2 + 1 + len(text))
	buf := make([]byte, l)

	binary.BigEndian.PutUint16(buf[0:2], l-2)
	buf[2] = '+'
	copy(buf[3:], []byte(text))

	log.Print(buf)

	binary.Write(conn, binary.BigEndian, &buf)
}

func handleDebugPacket(packet []byte) {
	text := packet[3:]
	log.Printf("[DEBUG PACKET] %s", text)
}

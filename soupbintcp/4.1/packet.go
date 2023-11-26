package soupbintcp

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	PacketTypeLoginRequest   = 'L'
	PacketLengthLoginRequest = 52

	PacketTypeLoginAccepted   = 'A'
	PacketLengthLoginAccepted = 31

	PacketTypeLoginRejected   = 'J'
	PacketLengthLoginRejected = 2

	PacketTypeLogoutRequest   = 'O'
	PacketLengthLogoutRequest = 1

	PacketTypeClientHeartbeat   = 'R'
	PacketLengthClientHeartbeat = 1

	PacketTypeServerHeartbeat   = 'H'
	PacketLengthServerHeartbeat = 1

	PacketTypeEndOfSession   = 'Z'
	PacketLengthEndOfSession = 1

	// Variable length packets
	PacketTypeSequencedData   = 'S'
	PacketTypeUnsequencedData = 'U'
	PacketTypeDebug           = '+'

	LoginRejectedNotAuthorized      = 'A'
	LoginRejectedSessionUnavailable = 'S'
)

type Packet interface {
	Bytes() []byte
}

func getNextPacket(conn net.Conn) ([]byte, error) {
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

type Header struct {
	Length [2]byte
	Type   byte
}

func (h Header) Bytes() []byte {
	buf := make([]byte, 3)

	copy(buf[0:], h.Length[:])
	buf[2] = h.Type

	return buf
}

type HeartbeatPacket struct {
	Header
}

func makeClientHeartbeatPacket() Packet {
	return HeartbeatPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthClientHeartbeat},
			Type:   PacketTypeClientHeartbeat,
		},
	}
}

func makeServerHeartbeatPacket() Packet {
	return HeartbeatPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthServerHeartbeat},
			Type:   PacketTypeServerHeartbeat,
		},
	}
}

type LoginAcceptedPacket struct {
	Header
	Session        [10]byte
	SequenceNumber [20]byte
}

func makeLoginAcceptedPacket(session string, sequence uint64) Packet {
	packet := LoginAcceptedPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthLoginAccepted},
			Type:   PacketTypeLoginAccepted,
		},
	}

	copy(packet.SequenceNumber[:], fmt.Sprintf("%20d", sequence))
	copy(packet.Session[:], []byte(session))

	return packet
}

func (p LoginAcceptedPacket) Bytes() []byte {
	buf := make([]byte, 33)

	copy(buf[0:], p.Header.Bytes())
	copy(buf[3:], p.Session[:])
	copy(buf[13:], p.SequenceNumber[:])

	return buf
}

type LoginRejectedPacket struct {
	Header
	Reason byte
}

func makeLoginRejectedPacket(reason byte) Packet {
	return LoginRejectedPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthLoginRejected},
			Type:   PacketTypeLoginRejected,
		},
		Reason: reason,
	}
}

func (p LoginRejectedPacket) Bytes() []byte {
	buf := make([]byte, 4)

	copy(buf[0:], p.Header.Bytes())
	buf[3] = p.Reason

	return buf
}

type LoginRequestPacket struct {
	Header
	Username         [6]byte
	Password         [10]byte
	Session          [10]byte
	SequenceNumber   [20]byte
	HeartbeatTimeout [5]byte
}

func makeLoginRequestPacket(username, password, session string, sequence uint64) Packet {
	packet := LoginRequestPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthLoginRequest},
			Type:   PacketTypeLoginRequest,
		},
	}

	username = fmt.Sprintf("%-6s", username)
	password = fmt.Sprintf("%-10s", password)

	copy(packet.Username[:], username)
	copy(packet.Password[:], password)
	copy(packet.HeartbeatTimeout[:], fmt.Sprint(heartbeatPeriod_ms))
	copy(packet.SequenceNumber[:], fmt.Sprintf("%20s", strconv.Itoa(int(sequence))))
	copy(packet.Session[:], fmt.Sprintf("%10s", session))

	return packet
}

func (p LoginRequestPacket) Bytes() []byte {
	buf := make([]byte, PacketLengthLoginRequest+2) // +2 for the packet length field

	copy(buf[0:], p.Header.Bytes())
	copy(buf[3:], p.Username[:])
	copy(buf[9:], p.Password[:])
	copy(buf[19:], p.Session[:])
	copy(buf[29:], p.SequenceNumber[:])
	copy(buf[49:], p.HeartbeatTimeout[:])

	return buf
}

type LogoutRequestPacket struct {
	Header
}

func makeLogoutRequestPacket() Packet {
	return LogoutRequestPacket{
		Header: Header{
			Length: [2]byte{0, PacketLengthLogoutRequest},
			Type:   PacketTypeLogoutRequest,
		},
	}
}

type SequencedDataPacket struct {
	Header
	Message []byte
}

func (p SequencedDataPacket) makeSequencedDataPacket(data []byte) Packet {
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(1+len(data))) // +1 is for the type field
	packetLength := ([2]byte)(l)

	return SequencedDataPacket{
		Header: Header{
			Length: packetLength,
			Type:   PacketTypeSequencedData,
		},
		Message: data,
	}
}

func (p SequencedDataPacket) Bytes() []byte {
	buf := make([]byte, 3+len(p.Message)) // +3 is for header

	copy(buf[0:], p.Header.Bytes())
	copy(buf[3:], []byte(p.Message))

	return buf
}

type UnsequencedDataPacket struct {
	Header
	Message []byte
}

func makeUnsequencedDataPacket(data []byte) Packet {
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(1+len(data))) // +1 is for the type field
	packetLength := ([2]byte)(l)

	return UnsequencedDataPacket{
		Header: Header{
			Length: packetLength,
			Type:   PacketTypeUnsequencedData,
		},
		Message: data,
	}
}

func (p UnsequencedDataPacket) Bytes() []byte {
	buf := make([]byte, 3+len(p.Message)) // +3 is for header

	copy(buf[0:], p.Header.Bytes())
	copy(buf[3:], []byte(p.Message))

	return buf
}

type DebugPacket struct {
	Header
	Text string
}

func makeDebugPacket(text string) Packet {
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(1+len(text))) // +1 is for the type field
	packetLength := ([2]byte)(l)

	return DebugPacket{
		Header: Header{
			Length: packetLength,
			Type:   PacketTypeUnsequencedData,
		},
		Text: text,
	}
}

func (p DebugPacket) Bytes() []byte {
	buf := make([]byte, 3+len(p.Text)) // +3 is for the header

	copy(buf[0:], p.Header.Bytes())
	copy(buf[3:], []byte(p.Text))

	return buf
}

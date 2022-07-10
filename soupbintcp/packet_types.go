package soupbintcp

type Packet struct {
	Length [2]byte
	Type   byte
}

type DebugPacket struct {
	Packet
	Text string
}

type HeartbeatPacket struct {
	Packet
}

type LoginAcceptedPacket struct {
	Packet
	Session        [10]byte
	SequenceNumber [20]byte
}

type LoginRejectedPacket struct {
	Packet
	Reason byte
}

type LoginRequestPacket struct {
	Packet
	Username         [6]byte
	Password         [10]byte
	Session          [10]byte
	SequenceNumber   [20]byte
	HeartbeatTimeout [5]byte
}

type LogoutRequestPacket struct {
	Packet
}

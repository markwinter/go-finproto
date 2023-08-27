package soupbintcp

const (
	PacketLoginRequest  = 'L'
	PacketLoginAccepted = 'A'
	PacketLoginRejected = 'J'
	PacketLogoutRequest = 'O'

	PacketSequencedData   = 'S'
	PacketUnsequencedData = 'U'

	PacketClientHeartbeat = 'R'
	PacketServerHeartbeat = 'H'

	PacketDebug        = '+'
	PacketEndOfSession = 'Z'

	LoginRejectedNotAuthorized      = 'A'
	LoginRejectedSessionUnavailable = 'S'
)

type Header struct {
	Length [2]byte
	Type   byte
}

type DebugPacket struct {
	Header
	Text string
}

type HeartbeatPacket struct {
	Header
}

type LoginAcceptedPacket struct {
	Header
	Session        [10]byte
	SequenceNumber [20]byte
}

type LoginRejectedPacket struct {
	Header
	Reason byte
}

type LoginRequestPacket struct {
	Header
	Username         [6]byte
	Password         [10]byte
	Session          [10]byte
	SequenceNumber   [20]byte
	HeartbeatTimeout [5]byte
}

type LogoutRequestPacket struct {
	Header
}

type SequencedDataPacket struct {
	Header
	Message []byte
}

type UnsequencedDataPacket struct {
	Header
	Message []byte
}

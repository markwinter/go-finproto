package soupbintcp

import (
	"encoding/binary"
	"log"
	"net"
)

type session struct {
	Id                    string
	CurrentSequenceNumber uint64

	conns     []net.Conn
	dataStore dataStore
}

func makeSession(id string) session {
	s := session{
		Id:                    id,
		CurrentSequenceNumber: 1,
		dataStore:             &memoryStore{},
	}
	return s
}

func (s *session) addConn(conn net.Conn) {
	s.conns = append(s.conns, conn)
}

func (s *session) send(data []byte) error {
	if err := s.dataStore.store(data); err != nil {
		return err
	}

	// TODO: make this better
	for _, conn := range s.conns {
		if err := binary.Write(conn, binary.BigEndian, &data); err != nil {
			log.Printf("failed sending to %q\n", conn.RemoteAddr())
		}
	}

	return nil
}

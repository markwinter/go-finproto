package soupbintcp

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
)

type dataStore interface {
	Store(data []byte) error
	Read(index int) ([]byte, error)
}

type memoryStore struct {
	data [][]byte
}

func (s *memoryStore) Store(data []byte) error {
	s.data = append(s.data, data)
	return nil
}
func (s *memoryStore) Read(index int) ([]byte, error) {
	if index > len(s.data) {
		return []byte{}, errors.New("index past data store length")
	}
	return s.data[index], nil
}

type Session struct {
	Id                    string
	CurrentSequenceNumber uint64

	conns     []net.Conn
	dataStore dataStore
}

func MakeSession(id string) Session {
	s := Session{
		Id:                    id,
		CurrentSequenceNumber: 1,
		dataStore:             &memoryStore{},
	}
	return s
}

func (s *Session) AddConn(conn net.Conn) {
	s.conns = append(s.conns, conn)
}

func (s *Session) Send(data []byte) error {
	if err := s.dataStore.Store(data); err != nil {
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

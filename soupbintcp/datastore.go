package soupbintcp

import "errors"

type dataStore interface {
	store(data []byte) error
	read(index int) ([]byte, error)
	clear()
}

type memoryStore struct {
	data [][]byte
}

func (s *memoryStore) store(data []byte) error {
	s.data = append(s.data, data)
	return nil
}

func (s *memoryStore) read(index int) ([]byte, error) {
	if index > len(s.data) {
		return []byte{}, errors.New("index past data store length")
	}
	return s.data[index], nil
}

func (s *memoryStore) clear() {
	s.data = nil
}

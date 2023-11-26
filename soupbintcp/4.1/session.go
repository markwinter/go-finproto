package soupbintcp

type session struct {
	id                    string
	currentSequenceNumber uint64
	dataStore             dataStore
}

func makeSession(id string) session {
	s := session{
		id:                    id,
		currentSequenceNumber: 1,
		dataStore:             &memoryStore{},
	}
	return s
}

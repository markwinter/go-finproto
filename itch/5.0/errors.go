package itch

import "fmt"

type ErrInvalidPacketSize struct {
	err error
}

func (e ErrInvalidPacketSize) Error() string {
	return e.err.Error()
}

func NewInvalidPacketSize(w, g int) ErrInvalidPacketSize {
	return ErrInvalidPacketSize{
		err: fmt.Errorf("expected data len=%d but got=%d", w, g),
	}
}

type ErrInvalidPacketType struct {
	err error
}

func (e ErrInvalidPacketType) Error() string {
	return e.err.Error()
}

func NewInvalidPacketType(t byte) ErrInvalidPacketType {
	return ErrInvalidPacketType{
		err: fmt.Errorf("invalid packet type=%v", t),
	}
}

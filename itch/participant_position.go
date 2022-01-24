/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type MMMode uint8
type MMState uint8

var (
	MarketParticipants = make(map[string][]ParticipantPosition)
)

const (
	MMMODE_NORMAL        MMMode = 'N'
	MMMODE_PASSIVE       MMMode = 'P'
	MMMODE_SYNDICATE     MMMode = 'S'
	MMMODE_PRE_SYNDICATE MMMode = 'R'
	MMMODE_PENALTY       MMMode = 'L'

	MMSTATE_ACTIVE    MMState = 'A'
	MMSTATE_EXCUSED   MMState = 'E'
	MMSTATE_WITHDRAWN MMState = 'W'
	MMSTATE_SUSPENDED MMState = 'S'
	MMSTATE_DELETED   MMState = 'D'
)

type ParticipantPosition struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Mpid           string
	Stock          string
	PrimaryMM      bool
	Mode           MMMode
	State          MMState
}

func MakeParticipantPosition(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	primary := false
	if data[23] == 'Y' {
		primary = true
	}

	pp := ParticipantPosition{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Mpid:           strings.TrimSpace(string(data[11:15])),
		Stock:          strings.TrimSpace(string(data[15:23])),
		PrimaryMM:      primary,
		Mode:           MMMode(data[24]),
		State:          MMState(data[25]),
	}

	MarketParticipants[pp.Mpid] = append(MarketParticipants[pp.Mpid], pp)

	return pp
}

func (p *ParticipantPosition) String() string {
	return fmt.Sprintf("[Market Participant Position]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"MPID: %v\n"+
		"Stock: %v\n"+
		"Primary: %v\n"+
		"Mode: %v\n"+
		"State: %v\n",
		p.StockLocate, p.TrackingNumber, p.Timestamp,
		p.Mpid, p.Stock, p.PrimaryMM, p.Mode, p.State,
	)
}

func (m MMMode) String() string {
	switch m {
	case MMMODE_NORMAL:
		return "Normal"
	case MMMODE_PASSIVE:
		return "Passive"
	case MMMODE_SYNDICATE:
		return "Syndicate"
	case MMMODE_PRE_SYNDICATE:
		return "Pre-Syndicate"
	case MMMODE_PENALTY:
		return "Penalty"
	}

	return "Unkown MMMode"
}

func (m MMState) String() string {
	switch m {
	case MMSTATE_ACTIVE:
		return "Active"
	case MMSTATE_EXCUSED:
		return "Excused/Withdrawn"
	case MMSTATE_WITHDRAWN:
		return "Withdrawn"
	case MMSTATE_SUSPENDED:
		return "Suspended"
	case MMSTATE_DELETED:
		return "Deleted"
	}

	return "Unknown MMState"
}

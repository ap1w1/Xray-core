package hysteria

import "github.com/xtls/xray-core/common/errors"

const maxUDPFragmentCount = 255

func FragUDPMessage(m *UDPMessage, maxSize int) ([]UDPMessage, error) {
	if maxSize <= 0 {
		return nil, errors.New("invalid max UDP message size")
	}
	if m.Size() <= maxSize {
		return []UDPMessage{*m}, nil
	}
	fullPayload := m.Data
	maxPayloadSize := maxSize - m.HeaderSize()
	if maxPayloadSize <= 0 {
		return nil, errors.New("max UDP message size is smaller than hysteria header")
	}
	fragCount := (len(fullPayload) + maxPayloadSize - 1) / maxPayloadSize // round up
	if fragCount > maxUDPFragmentCount {
		return nil, errors.New("UDP message requires too many fragments: ", fragCount)
	}
	frags := make([]UDPMessage, fragCount)
	for i, off := 0, 0; off < len(fullPayload); i++ {
		payloadSize := len(fullPayload) - off
		if payloadSize > maxPayloadSize {
			payloadSize = maxPayloadSize
		}
		frag := *m
		frag.FragID = uint8(i)
		frag.FragCount = uint8(fragCount)
		frag.Data = fullPayload[off : off+payloadSize]
		frags[i] = frag
		off += payloadSize
	}
	return frags, nil
}

type defragState struct {
	frags []*UDPMessage
	count uint8
	size  int
}

// Defragger handles defragmentation of interleaved UDP messages.
type Defragger struct {
	packets map[uint16]*defragState
}

func (d *Defragger) Feed(m *UDPMessage) *UDPMessage {
	if m.FragCount <= 1 {
		return m
	}
	if m.PacketID == 0 || m.FragID >= m.FragCount {
		return nil
	}
	if d.packets == nil {
		d.packets = make(map[uint16]*defragState)
	}

	state := d.packets[m.PacketID]
	if state == nil || len(state.frags) != int(m.FragCount) {
		state = &defragState{frags: make([]*UDPMessage, m.FragCount)}
		d.packets[m.PacketID] = state
	}
	if state.frags[m.FragID] != nil {
		return nil
	}

	state.frags[m.FragID] = m
	state.count++
	state.size += len(m.Data)
	if int(state.count) != len(state.frags) {
		return nil
	}

	data := make([]byte, state.size)
	off := 0
	for _, frag := range state.frags {
		if frag == nil {
			return nil
		}
		off += copy(data[off:], frag.Data)
	}
	delete(d.packets, m.PacketID)

	assembled := *m
	assembled.Data = data
	assembled.FragID = 0
	assembled.FragCount = 1
	return &assembled
}

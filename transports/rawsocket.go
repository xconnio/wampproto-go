package transports

import (
	"errors"
	"fmt"
	"math"
)

const (
	MAGIC              byte = 0x7F
	ProtocolMaxMsgSize      = 1 << 24
	DefaultMaxMsgSize       = 1 << 20

	SerializerJson    Serializer = 1
	SerializerMsgpack Serializer = 2
	SerializerCbor    Serializer = 3

	MessageWamp Message = 0
	MessagePing Message = 1
	MessagePong Message = 2
)

type Serializer int
type Message byte

type Handshake struct {
	serializer     Serializer
	maxMessageSize int
}

func NewHandshake(serializer Serializer, maxMessageSize int) *Handshake {
	return &Handshake{
		serializer:     serializer,
		maxMessageSize: maxMessageSize,
	}
}

func (h *Handshake) Serializer() Serializer {
	return h.serializer
}

func (h *Handshake) MaxMessageSize() int {
	return h.maxMessageSize
}

func SendHandshake(hs *Handshake) ([]byte, error) {
	if hs.MaxMessageSize() > ProtocolMaxMsgSize {
		return nil, errors.New("maxMessageSize must not be more than 16 megabytes")
	}

	log2 := int(math.Log2(float64(hs.MaxMessageSize())))
	if (1<<log2) != hs.MaxMessageSize() || log2 < 9 {
		return nil, errors.New("maxMessageSize must be a power of 2 and >= 512")
	}

	b1 := byte((log2-9)<<4 | (int(hs.Serializer()) & 0x0F))
	return []byte{MAGIC, b1, 0x00, 0x00}, nil
}

func ReceiveHandshake(data []byte) (*Handshake, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("expected 4 bytes for handshake response, got %d", len(data))
	}
	if data[0] != MAGIC {
		return nil, fmt.Errorf("expected MAGIC, got %d", data[0])
	}
	if data[2] != 0x00 || data[3] != 0x00 {
		return nil, fmt.Errorf("expected 0x00 for third and fourth byte, got %d and %d", data[2], data[3])
	}

	serializer := data[1] & 0x0F
	sizeShift := (data[1] >> 4) + 9
	maxMessageSize := 1 << sizeShift

	return NewHandshake(Serializer(serializer), maxMessageSize), nil
}

type MessageHeader struct {
	kind   Message
	length int
}

func NewMessageHeader(kind Message, length int) *MessageHeader {
	return &MessageHeader{
		kind:   kind,
		length: length,
	}
}

func (h *MessageHeader) Kind() Message {
	return h.kind
}

func (h *MessageHeader) Length() int {
	return h.length
}

func SendMessageHeader(header *MessageHeader) []byte {
	data := IntToBytes(header.Length())
	return []byte{
		byte(header.Kind()),
		data[0],
		data[1],
		data[2],
	}
}

func ReceiveMessageHeader(data []byte) (*MessageHeader, error) {
	return NewMessageHeader(Message(data[0]), BytesToInt(data[1:])), nil
}

func IntToBytes(i int) []byte {
	return []byte{
		byte((i >> 16) & 0xFF),
		byte((i >> 8) & 0xFF),
		byte(i & 0xFF),
	}
}

func BytesToInt(b []byte) int {
	n := 0
	for _, byteVal := range b {
		n = (n << 8) | int(byteVal)
	}

	return n
}

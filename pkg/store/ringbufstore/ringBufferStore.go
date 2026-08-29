package rbs

import (
	"github.com/Joey574/sink/v2/pkg/ds/rb"
	"github.com/Joey574/sink/v2/pkg/sink"
)

type RingBufferStore struct {
	rb    *rb.RingBuffer
	chunk []byte
}

func New(capacity int) *RingBufferStore {
	return &RingBufferStore{
		rb:    rb.New(capacity),
		chunk: make([]byte, min(capacity, 4096)),
	}
}

func (s *RingBufferStore) Write(level sink.LogLevel, p []byte) error {
	data := append([]byte{0x00, 0xCC, byte(level)}, p...)
	data = append(data, []byte{0xCC, 0x00}...)

	_, err := s.rb.Write(data)
	return err
}

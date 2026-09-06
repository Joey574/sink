package rb

import (
	"fmt"
	"io"
	"sync"
)

type RingBuffer struct {
	capacity int
	size     int
	start    int
	end      int
	version  uint64
	buf      []byte

	mx sync.Mutex
}

func New(capacity int) *RingBuffer {
	if capacity == 0 {
		panic("size must be greater than 0")
	}

	return &RingBuffer{
		capacity: capacity,
		size:     0,
		start:    0,
		end:      0,
		version:  0,
		buf:      make([]byte, capacity),
	}
}

func NewBuffer(buf []byte) *RingBuffer {
	if len(buf) == 0 {
		panic("size must be greater than 0")
	}

	return &RingBuffer{
		capacity: len(buf),
		size:     0,
		start:    0,
		end:      0,
		version:  0,
		buf:      buf,
	}
}

func (r *RingBuffer) Close() error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.buf = nil
	r.capacity = 0
	r.size = 0
	r.start = 0
	r.end = 0
	r.version = 0
	return nil
}

func (r *RingBuffer) Version() uint64 {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.version
}

func (r *RingBuffer) Capacity() int {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.capacity
}

func (r *RingBuffer) Read(p []byte) (int, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.buf == nil {
		return 0, fmt.Errorf("buffer is closed")
	}

	n, err := r.readBufferLocked(p)
	if err == nil {
		r.start = (r.start + n) % r.capacity
		r.size -= n
	}
	return n, err
}

func (r *RingBuffer) Peek(p []byte) (int, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.buf == nil {
		return 0, fmt.Errorf("buffer is closed")
	}

	return r.readBufferLocked(p)
}

func (r *RingBuffer) Write(p []byte) (int, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.buf == nil {
		return 0, fmt.Errorf("buffer is closed")
	}

	n, err := r.writeLocked(p)
	if err == nil {
		r.version++
	}

	return n, err
}

func (r *RingBuffer) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

func (r *RingBuffer) PeekAll() ([]byte, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.buf == nil {
		return nil, fmt.Errorf("buffer is closed")
	}

	buf := make([]byte, r.size)
	_, err := r.readBufferLocked(buf)
	return buf, err
}

func (r *RingBuffer) ReadAll() ([]byte, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.buf == nil {
		return nil, fmt.Errorf("buffer is closed")
	}

	buf := make([]byte, r.size)
	_, err := r.readBufferLocked(buf)
	if err == nil {
		r.start = 0
		r.end = 0
		r.size = 0
	}
	return buf, err
}

func (r *RingBuffer) IsFull() bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.isFullLocked()
}

func (r *RingBuffer) IsEmpty() bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	return r.isEmptyLocked()
}

func (r *RingBuffer) ReadFrom(reader io.Reader) (int64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	const chunkSize = 4096
	var err error

	chunk := make([]byte, chunkSize)
	n := int64(0)

	for err != io.EOF {
		count, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return n, err
		}

		n += int64(count)
		_, werr := r.writeLocked(chunk[:count])
		if werr != nil {
			return n, err
		}
	}

	return n, nil
}

func (r *RingBuffer) WriteTo(w io.Writer) (int64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	const chunkSize = 4096
	var err error

	idx := min(r.capacity, r.capacity-r.start)
	n, err := w.Write(r.buf[r.start:idx])
	if err != nil {
		return int64(n), err
	}

	if idx != r.capacity {
		n2, err := w.Write(r.buf[:r.end])
		n += n2

		if err != nil {
			return int64(n), err
		}
	}

	return int64(n), nil
}

func (r *RingBuffer) Find(needle, p []byte) int {
	r.mx.Lock()
	defer r.mx.Unlock()

	idx := r.start
	nidx := 0
	for range r.size {
		if r.buf[idx] == needle[nidx] {
			nidx++
		}

		if nidx == len(needle) {
			// needle found, read from idx-nidx into p
			r.readBufferLocked(p)
		}

		idx++
	}

	return -1
}

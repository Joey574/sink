package rb

import "io"

func (r *RingBuffer) writeLocked(p []byte) (int, error) {
	if len(p) >= r.capacity {
		return r.truncateCopyLocked(p)
	}

	if overflow := r.size + len(p) - r.capacity; overflow > 0 {
		r.start = (r.start + overflow) % r.capacity
	}

	r.size = min(r.size+len(p), r.capacity)
	n1 := copy(r.buf[r.end:], p)
	copy(r.buf, p[n1:])

	r.end = (r.end + len(p)) % r.capacity
	return len(p), nil
}

func (r *RingBuffer) isFullLocked() bool {
	return r.start == r.end && r.size != 0
}

func (r *RingBuffer) isEmptyLocked() bool {
	return r.start == r.end && r.size == 0
}

// Handles copy for data which is larger than the buffer, requiring truncation
func (r *RingBuffer) truncateCopyLocked(p []byte) (int, error) {
	idx := len(p) - len(r.buf)
	copy(r.buf, p[idx:])
	r.start = 0
	r.end = 0
	r.size = r.capacity
	return len(r.buf), nil
}

func (r *RingBuffer) readBufferLocked(p []byte) (int, error) {
	n := min(len(p), r.size)
	if n == 0 {
		return 0, io.EOF
	}

	idx := min(n, r.capacity-r.start)
	copy(p[:idx], r.buf[r.start:r.start+idx])

	if idx < n {
		copy(p[idx:n], r.buf[:n-idx])
	}

	return n, nil
}

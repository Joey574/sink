package memoryStore

import (
	"sync"
	"time"

	"github.com/Joey574/sink/v2/pkg/sink"
)

type Log struct {
	level sink.LogLevel
	time  time.Time
	data  []byte
}

type MemoryStore struct {
	mx      sync.Mutex
	logs    []Log
	version uint64
}

func New() *MemoryStore {
	return &MemoryStore{
		logs: make([]Log, 0),
	}
}

func (ms *MemoryStore) Write(level sink.LogLevel, p []byte) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	log := Log{
		level: level,
		time:  time.Now(),
	}
	copy(log.data, p)

	ms.version++
	ms.logs = append(ms.logs, log)
	return nil
}

func (ms *MemoryStore) Version() uint64 {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	return ms.version
}

func (ms *MemoryStore) Pop() Log {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	log := ms.logs[len(ms.logs)-1]
	ms.logs = ms.logs[:len(ms.logs)-1]
	return log
}

func (ms *MemoryStore) Peek() Log {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	log := ms.logs[len(ms.logs)-1]
	return log
}

func (ms *MemoryStore) Logs() []Log {
	return ms.logs
}

package sink

import (
	"io"
	"sync"
)

type threadSafeSink struct {
	snk Sink
	mx  sync.Mutex
}

// Increments log level by 1 up to max of TRACE
func (s *threadSafeSink) IncLogLevel() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.snk.IncLogLevel()
}

// Decrements log level by 1 to a min of QUIET
func (s *threadSafeSink) DecLogLevel() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.snk.DecLogLevel()
}

// Sets log level directly
func (s *threadSafeSink) SetLogLevel(l LogLevel) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.SetLogLevel(l)
}

// Appends writers to the set of sinks
func (s *threadSafeSink) PushSinks(w ...io.Writer) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.PushSinks(w...)
}

// Append stores
func (s *threadSafeSink) PushStores(w ...Store) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.PushStores(w...)
}

// Pops the first sinks and returns it
func (s *threadSafeSink) PopSink() io.Writer {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.PopSink()
}

// Pops the first store and returns it
func (s *threadSafeSink) PopStore() Store {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.PopStore()
}

// Sets sinks to nil
func (s *threadSafeSink) FlushSinks() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.FlushSinks()
}

// Sets stores to nil
func (s *threadSafeSink) FlushStores() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.FlushStores()
}

// Sets the logging format string to be used, by default this is empty
// \d => outputs the datetime
// \t => writes the log level as a string
// \c => writes the caller name
// * => log data
func (s *threadSafeSink) SetFormat(f string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.SetFormat(f)
}

// Writes the set of bytes to the sink, exits early in event of error and returns it
func (s *threadSafeSink) Write(l LogLevel, p []byte) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.Write(l, p)
}

func (s *threadSafeSink) WriteString(l LogLevel, w string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.WriteString(l, w)
}

func (s *threadSafeSink) Print(l LogLevel, a ...any) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.Print(l, a...)
}

func (s *threadSafeSink) Println(l LogLevel, a ...any) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.Println(l, a...)
}

func (s *threadSafeSink) Printf(l LogLevel, format string, a ...any) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.snk.Printf(l, format, a...)
}

func (s *threadSafeSink) Fatal(v ...any) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.snk.Fatal(v...)
}

func (s *threadSafeSink) Fatalln(v ...any) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.snk.Fatalln(v...)
}

func (s *threadSafeSink) Fatalf(format string, a ...any) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.snk.Fatalf(format, a...)
}

package sink

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type threadUnsafeSink struct {
	level  LogLevel
	format string
	sinks  []io.Writer
	stores []Store

	parent Sink
	name   string
}

func (s *threadUnsafeSink) SetParent(p Sink) {
	s.parent = p
}

func (s *threadUnsafeSink) Parent() Sink {
	return s.parent
}

func (s *threadUnsafeSink) SetName(n string) {
	s.name = n
}

func (s *threadUnsafeSink) Name() string {
	return s.name
}

func (s *threadUnsafeSink) CallStack() string {
	if s.parent == nil {
		return s.name
	}

	return fmt.Sprintf("%s.%s", s.parent.CallStack(), s.name)
}

// Increments log level by 1 up to max of TRACE
func (s *threadUnsafeSink) IncLogLevel() {
	s.level = min(TRACE, s.level+1)
}

// Decrements log level by 1 to a min of QUIET
func (s *threadUnsafeSink) DecLogLevel() {
	s.level = max(QUIET, s.level-1)
}

// Sets log level directly
func (s *threadUnsafeSink) SetLogLevel(l LogLevel) {
	s.level = l
}

// Appends writers to the set of sinks
func (s *threadUnsafeSink) PushSinks(w ...io.Writer) {
	s.sinks = append(s.sinks, w...)
}

// Append stores
func (s *threadUnsafeSink) PushStores(w ...Store) {
	s.stores = append(s.stores, w...)
}

// Pops the first sinks and returns it
func (s *threadUnsafeSink) PopSink() io.Writer {
	if len(s.sinks) == 0 {
		return nil
	}

	snk := s.sinks[0]
	s.sinks = s.sinks[1:]
	return snk
}

// Pops the first store and returns it
func (s *threadUnsafeSink) PopStore() Store {
	if len(s.stores) == 0 {
		return nil
	}

	st := s.stores[0]
	s.stores = s.stores[1:]
	return st
}

// Sets sinks to nil
func (s *threadUnsafeSink) FlushSinks() {
	s.sinks = nil
}

// Sets stores to nil
func (s *threadUnsafeSink) FlushStores() {
	s.stores = nil
}

// Sets the logging format string to be used, by default this is empty
// \d => outputs the datetime
// \t => writes the log level as a string
// \c => writes the caller name
// \s => writes the sink call stack
// * => log data
func (s *threadUnsafeSink) SetFormat(f string) error {
	s.format = f
	return nil
}

// Writes the set of bytes to the sink, exits early in event of error and returns it
func (s *threadUnsafeSink) Write(l LogLevel, p []byte) error {
	return s.writeLocked(l, p)
}

func (s *threadUnsafeSink) WriteString(l LogLevel, w string) error {
	return s.writeLocked(l, []byte(w))
}

func (s *threadUnsafeSink) Print(l LogLevel, a ...any) error {
	return s.writeLocked(l, fmt.Append(nil, a...))
}

func (s *threadUnsafeSink) Println(l LogLevel, a ...any) error {
	return s.writeLocked(l, fmt.Appendln(nil, a...))
}

func (s *threadUnsafeSink) Printf(l LogLevel, format string, a ...any) error {
	return s.writeLocked(l, fmt.Appendf(nil, format, a...))
}

func (s *threadUnsafeSink) Fatal(v ...any) {
	s.writeLocked(QUIET, fmt.Append(nil, v...))
	os.Exit(1)
}

func (s *threadUnsafeSink) Fatalln(v ...any) {
	s.writeLocked(QUIET, fmt.Appendln(nil, v...))
	os.Exit(1)
}

func (s *threadUnsafeSink) Fatalf(format string, a ...any) {
	s.writeLocked(QUIET, fmt.Appendf(nil, format, a...))
	os.Exit(1)
}

func (s *threadUnsafeSink) writeLocked(l LogLevel, p []byte) error {
	_ = s.writeToStores(l, p)

	if l > s.level {
		return nil
	}

	data := s.formatString(string(p), l)
	return s.writeToSinks([]byte(data))
}

func (s *threadUnsafeSink) writeToSinks(b []byte) error {
	for _, w := range s.sinks {
		if _, err := w.Write(b); err != nil {
			return err
		}
	}

	return nil
}

func (s *threadUnsafeSink) writeToStores(level LogLevel, p []byte) error {
	for _, store := range s.stores {
		if err := store.Write(level, p); err != nil {
			return err
		}
	}

	return nil
}

func (s *threadUnsafeSink) formatString(data string, level LogLevel) string {
	if s.format == "" {
		return data
	}

	out := strings.ReplaceAll(s.format, "*", data)
	out = strings.ReplaceAll(out, `\d`, time.Now().Format("2006-01-02 15:04:05"))
	out = strings.ReplaceAll(out, `\c`, extCallerName())
	out = strings.ReplaceAll(out, `\s`, s.CallStack())
	return strings.ReplaceAll(out, `\t`, levelString(level))
}

package sink

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogLevel int8

const (
	QUIET = LogLevel(iota) + 1
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

type Store interface {
	WriteWithCtx(context.Context, LogLevel, []byte) error
	Write(LogLevel, []byte) error
}

type Sink struct {
	level  LogLevel
	format string
	sinks  []io.Writer
	stores []Store
}

func New(options ...func(s *Sink)) *Sink {
	s := &Sink{
		level:  INFO,
		format: "",
	}

	for _, o := range options {
		o(s)
	}

	return s
}

// Increments log level by 1 up to max of TRACE
func (s *Sink) IncLogLevel() {
	s.level = min(TRACE, s.level+1)
}

// Decrements log level by 1 to a min of QUIET
func (s *Sink) DecLogLevel() {
	s.level = max(QUIET, s.level-1)
}

// Sets log level directly
func (s *Sink) SetLogLevel(l LogLevel) {
	s.level = l
}

// Appends writers to the set of sinks
func (s *Sink) PushSinks(w ...io.Writer) {
	s.sinks = append(s.sinks, w...)
}

// Pops the first sinks and returns it
func (s *Sink) PopSinks() io.Writer {
	if len(s.sinks) == 0 {
		return nil
	}

	snk := s.sinks[0]
	s.sinks = s.sinks[1:]
	return snk
}

// Sets sinks to nil
func (s *Sink) FlushSinks() {
	s.sinks = nil
}

// Sets the logging format string to be used, by default this is empty
// \d => outputs the datetime
// \t => writes the log level as a string
// \c => writes the caller name
// * => log data
func (s *Sink) SetFormat(f string) error {
	s.format = f
	return nil
}

// Writes the set of bytes to the sink, exits early in event of error and returns it
func (s *Sink) Write(l LogLevel, p []byte) error {
	return s.writeInternal(l, p, 1)
}

func (s *Sink) WriteString(l LogLevel, w string) error {
	return s.writeInternal(l, []byte(w), 1)
}

func (s *Sink) Printf(l LogLevel, format string, a ...any) error {
	return s.writeInternal(l, fmt.Appendf(nil, format, a...), 1)
}

func (s *Sink) Println(l LogLevel, a ...any) error {
	return s.writeInternal(l, fmt.Appendln(nil, a...), 1)
}

func (s *Sink) Print(l LogLevel, a ...any) error {
	return s.writeInternal(l, fmt.Append(nil, a...), 1)
}

func (s *Sink) Fatal(v ...any) {
	s.writeInternal(QUIET, fmt.Append(nil, v...), 1)
	os.Exit(1)
}

func (s *Sink) Fatalf(format string, a ...any) {
	s.writeInternal(QUIET, fmt.Appendf(nil, format, a...), 1)
	os.Exit(1)
}

func (s *Sink) Fatalln(v ...any) {
	s.writeInternal(QUIET, fmt.Appendln(nil, v...), 1)
	os.Exit(1)
}

func (s *Sink) writeInternal(l LogLevel, p []byte, depth int) error {
	_ = s.writeToStores(l, p)

	if l > s.level {
		return nil
	}

	data := s.formatString(string(p), l, depth)
	return s.writeToSinks([]byte(data))
}

func (s *Sink) writeToSinks(b []byte) error {
	for _, w := range s.sinks {
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sink) writeToStores(level LogLevel, p []byte) error {
	for _, s := range s.stores {
		err := s.Write(level, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sink) formatString(data string, level LogLevel, depth int) string {
	if s.format == "" {
		return data
	}

	out := strings.ReplaceAll(s.format, "*", data)
	out = strings.ReplaceAll(out, "\\d", time.Now().Format("2006-01-02 15:04:05"))

	out = strings.ReplaceAll(out, "\\c", callerName(2+depth))
	return strings.ReplaceAll(out, "\\t", levelString(level))
}

func (l LogLevel) String() string {
	return levelString(l)
}

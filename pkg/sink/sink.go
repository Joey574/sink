package sink

import (
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
// \t => writes the log type
// * => log data
// example
// [\d] [\t] *
func (s *Sink) SetFormat(f string) error {
	s.format = f
	return nil
}

// Writes the set of bytes to the sink, exits early in event of error and returns it
func (s *Sink) Write(l LogLevel, b []byte) error {
	_ = s.writeToStores(l, b)

	if l > s.level {
		return nil
	}

	data := s.formatString(string(b), l)
	return s.writeToSinks([]byte(data))
}

func (s *Sink) WriteString(l LogLevel, w string) error {
	return s.Write(l, []byte(w))
}

func (s *Sink) Printf(l LogLevel, format string, a ...any) error {
	return s.WriteString(l, fmt.Sprintf(format, a...))
}

func (s *Sink) Println(l LogLevel, a ...any) error {
	return s.WriteString(l, fmt.Sprintln(a...))
}

func (s *Sink) Print(l LogLevel, a ...any) error {
	return s.WriteString(l, fmt.Sprint(a...))
}

func (s *Sink) Fatal(v ...any) {
	s.Print(QUIET, v...)
	os.Exit(1)
}

func (s *Sink) Fatalf(format string, a ...any) {
	s.Printf(QUIET, format, a...)
	os.Exit(1)
}

func (s *Sink) Fatalln(v ...any) {
	s.Println(QUIET, v...)
	os.Exit(1)
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

func (s *Sink) formatString(data string, level LogLevel) string {
	if s.format == "" {
		return data
	}

	out := strings.ReplaceAll(s.format, "*", data)
	out = strings.ReplaceAll(out, "\\d", time.Now().Format("2006-01-02 15:04:05"))
	return strings.ReplaceAll(out, "\\t", levelString(level))
}

func (l LogLevel) String() string {
	return levelString(l)
}

func levelString(level LogLevel) string {
	switch level {
	case QUIET:
		// the only logs that output when quiet is selected should be fatal ones
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		return "UNKN"
	}
}

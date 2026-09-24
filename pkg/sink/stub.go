package sink

import (
	"fmt"
	"os"
)

type stub struct {
	name string
	Sink
}

func (s *stub) Clone() Sink {
	return &stub{
		name: s.name,
		Sink: s.Sink,
	}
}

func (s *stub) SetParent(p Sink) {
	s.Sink = p
}

func (s *stub) Parent() Sink {
	return s.Sink
}

func (s *stub) SetName(n string) {
	s.name = n
}

func (s *stub) Name() string {
	return s.name
}

func (s *stub) CallStack() string {
	return fmt.Sprintf("%s.%s", s.Sink.CallStack(), s.name)
}

func (s *stub) emit(l LogLevel, suffix string, p []byte) error {
	stack := s.name
	if suffix != "" {
		stack += "." + suffix
	}

	return s.Sink.emit(l, stack, p)
}

// Write logs the raw bytes p at level l under this stubs call stack.
func (s *stub) Write(l LogLevel, p []byte) error {
	return s.emit(l, "", p)
}

// WriteString logs w at level l under this stubs call stack.
func (s *stub) WriteString(l LogLevel, w string) error {
	return s.emit(l, "", []byte(w))
}

// Print formats like fmt.Print and logs under this stubs call stack.
func (s *stub) Print(l LogLevel, a ...any) error {
	return s.emit(l, "", fmt.Append(nil, a...))
}

// Println formats like fmt.Println and logs under this stubs call stack.
func (s *stub) Println(l LogLevel, a ...any) error {
	return s.emit(l, "", fmt.Appendln(nil, a...))
}

// Printf formats like fmt.Printf and logs under this stubs call stack.
func (s *stub) Printf(l LogLevel, format string, a ...any) error {
	return s.emit(l, "", fmt.Appendf(nil, format, a...))
}

// Fatal logs like Print at the QUIET level and then exits with status 1.
func (s *stub) Fatal(v ...any) {
	s.emit(QUIET, "", fmt.Append(nil, v...))
	os.Exit(1)
}

// Fatalln logs like Println at the QUIET level and then exits with status 1.
func (s *stub) Fatalln(v ...any) {
	s.emit(QUIET, "", fmt.Appendln(nil, v...))
	os.Exit(1)
}

// Fatalf logs like Printf at the QUIET level and then exits with status 1.
func (s *stub) Fatalf(format string, a ...any) {
	s.emit(QUIET, "", fmt.Appendf(nil, format, a...))
	os.Exit(1)
}

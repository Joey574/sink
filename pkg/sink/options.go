package sink

import (
	"io"
	"os"
)

func EnableStdOut() func(*Sink) {
	return func(s *Sink) {
		s.PushSinks(os.Stdout)
	}
}

func SetLogLevel(level LogLevel) func(*Sink) {
	return func(s *Sink) {
		s.SetLogLevel(level)
	}
}

func SetFormat(format string) func(*Sink) {
	return func(s *Sink) {
		if err := s.SetFormat(format); err != nil {
			panic(err)
		}
	}
}

func PushSinks(w ...io.Writer) func(*Sink) {
	return func(s *Sink) {
		s.PushSinks(w...)
	}
}

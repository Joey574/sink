package sink

import "os"

func EnableStdOut() func(*Sink) {
	return func(s *Sink) {
		s.PushSinks(os.Stdout)
	}
}

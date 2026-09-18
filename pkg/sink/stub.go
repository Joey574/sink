package sink

import "fmt"

type stub struct {
	name string
	Sink
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

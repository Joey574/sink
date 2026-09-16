package sink

func New(options ...func(Sink) Sink) Sink {
	var s Sink
	s = &threadUnsafeSink{
		level:  INFO,
		format: "",
	}

	for _, o := range options {
		s = o(s)
	}

	return s
}

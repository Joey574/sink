package sink

import (
	"io"
	"os"
)

func EnableStdOut() func(Sink) Sink {
	return func(s Sink) Sink {
		s.PushSinks(os.Stdout)
		return s
	}
}

func SetLogLevel(level LogLevel) func(Sink) Sink {
	return func(s Sink) Sink {
		s.SetLogLevel(level)
		return s
	}
}

func SetFormat(format string) func(Sink) Sink {
	return func(s Sink) Sink {
		if err := s.SetFormat(format); err != nil {
			panic(err)
		}
		return s
	}
}

func PushSinks(w ...io.Writer) func(Sink) Sink {
	return func(s Sink) Sink {
		s.PushSinks(w...)
		return s
	}
}

func PushStores(w ...Store) func(Sink) Sink {
	return func(s Sink) Sink {
		s.PushStores(w...)
		return s
	}
}

func ThreadSafe() func(Sink) Sink {
	return func(s Sink) Sink {
		return &threadSafeSink{
			snk: s,
		}
	}
}

func SetParent(p Sink) func(Sink) Sink {
	return func(s Sink) Sink {
		s.SetParent(p)
		return s
	}
}

func SetName(n string) func(Sink) Sink {
	return func(s Sink) Sink {
		s.SetName(n)
		return s
	}
}

// This argument is location dependent, it will overwrite all other options passed before it
func Clone(c Sink) func(Sink) Sink {
	return func(s Sink) Sink {
		s = c.Clone()
		return s
	}
}

func Wrap(c Sink) func(Sink) Sink {
	return func(_ Sink) Sink {
		return &stub{
			Sink: c,
		}
	}
}

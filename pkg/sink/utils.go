package sink

import (
	"runtime"
	"strings"
)

func extCallerName() string {
	pcs := make([]uintptr, 32)

	n := runtime.Callers(2, pcs)
	if n == 0 {
		return "unknown"
	}

	frames := runtime.CallersFrames(pcs[:n])
	var internalPkg string

	for {
		frame, more := frames.Next()
		if frame.Function == "" {
			if !more {
				break
			}
			continue
		}

		pkg := pkgName(frame.Function)

		if internalPkg == "" {
			internalPkg = pkg
		} else if pkg != internalPkg {
			return frame.Function
		}

		if !more {
			break
		}
	}

	return "unknown"
}

func pkgName(funcName string) string {
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}

	dot := strings.IndexByte(funcName[lastSlash:], '.')
	if dot < 0 {
		return funcName
	}

	return funcName[:lastSlash+dot]
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

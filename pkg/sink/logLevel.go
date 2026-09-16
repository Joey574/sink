package sink

type LogLevel int8

const (
	QUIET = LogLevel(iota) + 1
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

func (l LogLevel) String() string {
	return levelString(l)
}

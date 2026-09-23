package sink

import (
	"context"
	"io"
)

type Store interface {
	WriteWithCtx(context.Context, LogLevel, []byte) error
	Write(LogLevel, []byte) error
}

type Sink interface {
	Clone() Sink

	SetParent(Sink)
	Parent() Sink

	SetName(string)
	Name() string

	CallStack() string

	IncLogLevel()
	DecLogLevel()
	SetLogLevel(LogLevel)

	PushSinks(w ...io.Writer)
	PushStores(w ...Store)

	PopSink() io.Writer
	PopStore() Store

	FlushSinks()
	FlushStores()

	SetFormat(f string) error

	Write(LogLevel, []byte) error
	WriteString(LogLevel, string) error

	Print(LogLevel, ...any) error
	Println(LogLevel, ...any) error
	Printf(LogLevel, string, ...any) error

	Fatal(...any)
	Fatalln(...any)
	Fatalf(string, ...any)

	emit(LogLevel, string, []byte) error
}

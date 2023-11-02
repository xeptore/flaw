package flaw

import (
	"runtime"
)

// P is shorthand for Record.Payload type
type P map[string]any

type Record struct {
	Function string
	Payload  P
}

type StackTrace struct {
	Line     int
	File     string
	Function string
}

type Flaw struct {
	// Inner is the error string that was passed in during initialization.
	Inner      string
	Records    []Record
	StackTrace []StackTrace
}

// Error satisfies builtin error interface type. It returns the inner error string.
func (f *Flaw) Error() string {
	return f.Inner
}

func traces() []StackTrace {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	st := make([]StackTrace, 0, n)
	for {
		frame, ok := frames.Next()
		if !ok {
			break
		}
		st = append(st, StackTrace{
			Line:     frame.Line,
			File:     frame.File,
			Function: frame.Function,
		})
	}
	return st
}

func callerFunc() string {
	const depth = 2
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, ok := frames.Next()
	if !ok {
		return "<UNKNOWN>"
	}
	return frame.Function
}

func newFlawWithoutTrace(err error) *Flaw {
	return &Flaw{
		Records:    nil,
		Inner:      err.Error(),
		StackTrace: nil,
	}
}

// From creates a [Flaw] instance from an existing error. You can append contextual
// information to it using the [Flaw.Append] function immediately after instantiation,
// or by the caller function, after making sure the returned error is of type [Flaw]
// (using [errors.As]), It panics if err is nil.
func From(err error) *Flaw {
	if nil == err {
		panic("err can not be nil")
	}
	f := newFlawWithoutTrace(err)
	f.StackTrace = traces()
	return f
}

// Append appends contextual information to [Flaw] instance. It can be called immediately
// after instantiation using [From], or by the parent caller function, after making sure
// the returned error is of type [Flaw] (using [errors.As]). It panics if payload is nil.
func (f *Flaw) Append(payload P) *Flaw {
	if nil == payload {
		panic("payload cannot be nil")
	}
	f.Records = append(f.Records, Record{
		Function: callerFunc(),
		Payload:  payload,
	})
	return f
}

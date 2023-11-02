package flaw

import (
	"runtime"
)

// Record contains JSON serialized contextual information object, and a key
// than can be used for logging purposes.
type Record struct {
	Function string         `json:"function"`
	Payload  map[string]any `json:"payload"`
}

type StackTrace struct {
	Line     int
	File     string
	Function string
}

type Flaw struct {
	Inner      string       `json:"inner"`
	Records    []Record     `json:"records"`
	StackTrace []StackTrace `json:"stack_trace"`
}

// Error satisfies [error]. It returns JSON serialized array of [Flaw].Records.
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

// From creates a [Flaw] instance from an existing error. It appends contextual
// information to it, if it already contains a [Flaw] inside (checked using
// [errors.As]), or creates a new instance similar to [New] with message, and
// err.Error concatenated together. It panics if err is nil.
func From(err error) *Flaw {
	if nil == err {
		panic("err can not be nil")
	}
	f := newFlawWithoutTrace(err)
	f.StackTrace = traces()
	return f
}

func (f *Flaw) Append(payload map[string]any) *Flaw {
	f.Records = append(f.Records, Record{
		Function: callerFunc(),
		Payload:  payload,
	})
	return f
}

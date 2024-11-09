package flaw

import (
	"fmt"
	"runtime"
)

// P is shorthand for Record.Payload type
type P map[string]any

type Record struct {
	Function string
	Payload  P
}

type StackTrace struct {
	// Line is the file line number of the location in this frame.
	// For non-leaf frames, this will be the location of a call.
	// This may be zero, if not known.
	Line int
	// File is the file name of the location in this frame.
	// For non-leaf frames, this will be the location of a call.
	// This may be the empty string if not known.
	File string
	// Function is the package path-qualified function name of
	// this call frame. If non-empty, this string uniquely
	// identifies a single function in the program.
	// This may be the empty string if not known.
	Function string
}

type JoinedError struct {
	// Message is the result of joined error Error method call.
	Message string
	// CallerStackTrace is the error generator stack trace,
	// which in a very rare case can be nil. See [runtime.CallersFrames],
	// and [runtime.Callers] for more information on when this might happen.
	CallerStackTrace *StackTrace
	// TypeName is the type name of the error, which is the result of calling [fmt.Sprintf("%T", err)].
	TypeName string
	// SyntaxRepr is the string representation of the error, which is the result of calling [fmt.Sprintf("%+#v", err)].
	SyntaxRepr string
}

type Flaw struct {
	// Inner is the error string that was passed in during initialization.
	Inner string
	// InnerType is the type name of the error, which is the result of calling [fmt.Sprintf("%T", err)].
	InnerType string
	// InnerSyntaxRepr is the string representation of the error, which is the result of calling [fmt.Sprintf("%+#v", err)].
	InnerSyntaxRepr string
	// JoinedErrors is the list of optional (nil-able) errors
	// joined into the error while traversing the stack back,
	// usually in the same function initiated the error.
	JoinedErrors []JoinedError
	// Records contains contextual information in order the error traversed the stack up,
	// i.e., the first item in the slice is the first record attached to the error,
	// and the last item is the most recent attached record.
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
	if n == 0 {
		return nil
	}
	frames := runtime.CallersFrames(pcs[:n])
	st := make([]StackTrace, 0, n)
	for {
		frame, ok := frames.Next()
		st = append(st, StackTrace{
			Line:     frame.Line,
			File:     frame.File,
			Function: frame.Function,
		})
		if !ok {
			break
		}
	}
	return st
}

func joinTrace() *StackTrace {
	var pc [1]uintptr
	n := runtime.Callers(3, pc[:])
	if n == 0 {
		return nil
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &StackTrace{
		Line:     frame.Line,
		File:     frame.File,
		Function: frame.Function,
	}
}

func callerFunc() string {
	const depth = 2
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return frame.Function
}

func newFlawWithoutTrace(err error) *Flaw {
	return &Flaw{
		Records:         nil,
		Inner:           err.Error(),
		InnerType:       fmt.Sprintf("%T", err),
		InnerSyntaxRepr: fmt.Sprintf("%+#v", err),
		StackTrace:      nil,
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
// Elements in payloads will be merged into payload in order they are provided,
// thus duplicate keys will be overwritten.
func (f *Flaw) Append(payload P, payloads ...P) *Flaw {
	if nil == payload {
		panic("payload must not be nil")
	}
	merged := make(P, len(payload))
	for k, v := range payload {
		merged[k] = v
	}
	for _, p := range payloads {
		for k, v := range p {
			merged[k] = v
		}
	}
	f.Records = append(f.Records, Record{
		Function: callerFunc(),
		Payload:  merged,
	})
	return f
}

// Join joins the error to the flaw as a JoinedError item.
// Usually, used in failed defer calls where the deferred function
// fails with another error, and it is desired to capture
// the information of the error, and attach it to the original error.
func (f *Flaw) Join(err error) *Flaw {
	if nil == err {
		panic("err must not be nil")
	}
	f.JoinedErrors = append(f.JoinedErrors, JoinedError{
		Message:          err.Error(),
		CallerStackTrace: joinTrace(),
		TypeName:         fmt.Sprintf("%T", err),
		SyntaxRepr:       fmt.Sprintf("%+#v", err),
	})
	return f
}

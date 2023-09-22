package flaw

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/xeptore/flaw/v5/internal/encoder"
)

var (
	// NewDict can be used to initialize a JSON object of contextual information record with a key.
	NewDict = encoder.Dict
)

type (
	Dict = encoder.Record
)

// Record contains JSON serialized contextual information object, and a key
// than can be used for logging purposes.
type Record struct {
	Message string
	Payload []byte
}

type StackTrace struct {
	Line     int
	File     string
	Function string
}

type Flaw struct {
	// ID is a 36 characters URL-safe unique identifier for the instance.
	ID      string
	Records []Record
	Traces  []StackTrace
}

// Error satisfies [error]. It returns JSON serialized array of [Flaw].Records.
func (f *Flaw) Error() string {
	var builder strings.Builder
	builder.WriteByte('[')
	for i, r := range f.Records {
		if i != 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(`{"message":`)
		msg, err := json.Marshal(r.Message)
		if nil != err {
			return fmt.Errorf("failed to serialize record message to json: %v", err).Error()
		}
		builder.Write(msg)
		builder.WriteString(`,"payload":`)
		if r.Payload == nil {
			builder.WriteString("null")
		} else {
			builder.Write(r.Payload)
		}
		builder.WriteString(`}`)
	}
	builder.WriteByte(']')
	return builder.String()
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

func mustGenerateID() string {
	for i := 0; i < 20; i++ {
		id, err := gonanoid.New(36)
		if nil == err {
			return id
		}
	}
	var str string
	for len(str) < 36 {
		str += strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	}
	return str[:36]
}

func newFlawWithoutTrace(message string) *Flaw {
	id := mustGenerateID()
	return &Flaw{
		ID: id,
		Records: []Record{
			{
				Message: message,
			},
		},
		Traces: nil,
	}
}

// New creates a [Flaw] instance with a message, and contextual information
// record embedded into it.
func New(message string) *Flaw {
	f := newFlawWithoutTrace(message)
	f.Traces = traces()
	return f
}

// From creates a [Flaw] instance from an existing error. It appends contextual
// information to it, if it already contains a [Flaw] inside (checked using
// [errors.As]), or creates a new instance similar to [New] with message, and
// err.Error concatenated together. It panics if err is nil.
func From(err error, message string) *Flaw {
	if nil == err {
		panic("err can not be nil")
	}
	if flaw := new(Flaw); errors.As(err, &flaw) {
		flaw.Records = append(flaw.Records, Record{
			Message: message,
		})
		return flaw
	}
	f := newFlawWithoutTrace(message + ": " + err.Error())
	f.Traces = traces()
	return f
}

func (f *Flaw) With(rec *encoder.Record) *Flaw {
	f.Records[len(f.Records)-1].Payload = encoder.JSON(rec)
	return f
}

package flaw

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/xeptore/flaw/v2/internal/encoder"
)

var (
	// Dict can be used to initialize a JSON object of contextual information record with a key.
	Dict = encoder.Dict
)

// Record contains JSON serialized contextual information object, and a key
// than can be used for logging purposes.
type Record struct {
	Key     string
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
		builder.WriteString(`{"key":"` + r.Key + `","payload":`)
		builder.Write(r.Payload)
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
		id, err := gonanoid.ID(36)
		if nil != err {
			continue
		}
		return id
	}
	var str string
	for len(str) < 36 {
		str += strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	}
	return str[:36]
}

func newFlawWithoutTrace(message string, rec *encoder.Record) *Flaw {
	id := mustGenerateID()
	return &Flaw{
		ID: id,
		Records: []Record{
			{
				Key:     rec.Key,
				Payload: encoder.JSON(encoder.AppendErr(rec, message)),
			},
		},
		Traces: nil,
	}
}

// New creates a [Flaw] instance with a message, and contextual information
// record embedded into it.
func New(message string, rec *encoder.Record) *Flaw {
	f := newFlawWithoutTrace(message, rec)
	f.Traces = traces()
	return f
}

// From creates a [Flaw] instance from an existing error. It appends contextual
// information to it, if it already contains a [Flaw] inside (checked using
// [errors.As]), or creates a new instance similar to [New] with message, and
// err.Error concatenated together. It panics if err is nil.
func From(err error, message string, rec *encoder.Record) *Flaw {
	if nil == err {
		panic("err can not be nil")
	}
	if flaw := new(Flaw); errors.As(err, &flaw) {
		flaw.Records = append(flaw.Records, Record{
			Key:     rec.Key,
			Payload: encoder.JSON(encoder.AppendErr(rec, message)),
		})
		return flaw
	}
	f := newFlawWithoutTrace(message+": "+err.Error(), rec)
	f.Traces = traces()
	return f
}

package flaw

import (
	"errors"
	"strings"

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

type Flaw struct {
	Records []Record
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

// New creates a [Flaw] instance with a message, and contextual information
// record embedded into it.
func New(message string, rec *encoder.Record) *Flaw {
	return &Flaw{
		Records: []Record{
			{
				Key:     rec.Key,
				Payload: encoder.JSON(encoder.AppendErr(rec, message)),
			},
		},
	}
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
	return New(message+": "+err.Error(), rec)
}

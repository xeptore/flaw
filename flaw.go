package flaw

import (
	"errors"
	"strings"

	"github.com/xeptore/flaw/v2/internal/encoder"
)

var (
	Dict = encoder.Dict
)

type Record struct {
	Key     string
	Payload []byte
}

type Flaw struct {
	Records []Record
}

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

// From creates a [Flaw] instance from an existing error.
// It is designed to be used in the middle of the application layers where you'd want to attach more info to the error you receive if it is already
// a [Flaw] error, or you'd want to initialize new one with some contextual info.
// It also works if [Flaw] is wrapped inside another error, since it uses [errors.As] semantic to extract [Flaw] type error.
// Although it is not required, but I would highly recommend to include a message,
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
	return &Flaw{
		Records: []Record{
			{
				Key:     rec.Key,
				Payload: encoder.JSON(encoder.AppendErr(rec, message+": "+err.Error())),
			},
		},
	}
}

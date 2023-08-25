package flaw

import (
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

func New(rec *encoder.Record) *Flaw {
	return &Flaw{
		Records: []Record{
			{
				Key:     rec.Key,
				Payload: rec.JSON(),
			},
		},
	}
}

func From(err error, rec *encoder.Record) *Flaw {
	if flaw, ok := err.(*Flaw); ok {
		flaw.Records = append(flaw.Records, Record{
			Key:     rec.Key,
			Payload: rec.JSON(),
		})
		return flaw
	}
	record := Record{
		Key:     rec.Key,
		Payload: rec.Err("error", err).JSON(),
	}
	return &Flaw{
		Records: []Record{record},
	}
}

package flaw

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Dict struct {
	key    string
	values map[string]any
}

func Key(key string) *Dict {
	return &Dict{
		key:    key,
		values: make(map[string]any),
	}
}

func (d *Dict) Int(key string, value int) *Dict {
	d.values[key] = value
	return d
}

func (d *Dict) Str(key string, value string) *Dict {
	d.values[key] = value
	return d
}

func (d *Dict) json() []byte {
	b, err := json.Marshal(d.values)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "failed to marshal dict values to json: %v", err)
	}
	return b
}

type Record struct {
	Key     string
	Message string
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
		builder.WriteString(`{"message":"` + r.Message + `","key":"` + r.Key + `","payload":`)
		builder.Write(r.Payload)
		builder.WriteString(`}`)
	}
	builder.WriteByte(']')
	return builder.String()
}

func New(message string, dict *Dict) *Flaw {
	return &Flaw{
		Records: []Record{
			{
				Key:     dict.key,
				Message: message,
				Payload: dict.json(),
			},
		},
	}
}

func From(err error, message string, dict *Dict) *Flaw {
	record := Record{
		Key:     dict.key,
		Message: message,
		Payload: dict.json(),
	}
	if flaw, ok := err.(*Flaw); ok {
		flaw.Records = append(flaw.Records, record)
		return flaw
	}
	return &Flaw{
		Records: []Record{record},
	}
}

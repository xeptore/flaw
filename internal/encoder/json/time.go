package json

import (
	"time"
)

// AppendTime formats the input time with [time.RFC3339] format
// and appends the encoded string to the input byte slice.
func (e Encoder) AppendTime(dst []byte, t time.Time) []byte {
	return append(t.AppendFormat(append(dst, '"'), time.RFC3339), '"')
}

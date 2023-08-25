package encoder

import (
	"net"
	"time"

	"github.com/xeptore/flaw/v2/internal/encoder/json"
)

var (
	enc = json.Encoder{}
)

const (
	errorKey = "error"
)

// Record is the container for the contextual information JSON object that is
// attached to a key. You can use method chaining approach to assign key-value
// pairs to an instance of Record. It does not handle field deduplication,
// and creates duplicate keys. In this case, many consumers will take the last
// value, but this is not guaranteed; check yours if in doubt.
type Record struct {
	Key string
	buf []byte
}

func Dict(key string) *Record {
	return &Record{
		Key: key,
		buf: enc.AppendBeginMarker(nil),
	}
}

func JSON(r *Record) []byte {
	r.buf = enc.AppendEndMarker(r.buf)
	return r.buf
}

func AppendErr(r *Record, msg string) *Record {
	r.buf = enc.AppendString(enc.AppendKey(r.buf, errorKey), msg)
	return r
}

func (r *Record) Bool(key string, b bool) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendBool(enc.AppendKey(r.buf, key), b)
	return r
}

func (r *Record) Bools(key string, b []bool) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendBools(enc.AppendKey(r.buf, key), b)
	return r
}

func (r *Record) Dict(key string, dict *Record) *Record {
	if key == errorKey {
		return r
	}
	dict.buf = enc.AppendEndMarker(dict.buf)
	r.buf = append(enc.AppendKey(r.buf, key), dict.buf...)
	return r
}

func (r *Record) Float32(key string, f float32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendFloat32(enc.AppendKey(r.buf, key), f)
	return r
}

func (r *Record) Float64(key string, f float64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendFloat64(enc.AppendKey(r.buf, key), f)
	return r
}

func (r *Record) Floats32(key string, f []float32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendFloats32(enc.AppendKey(r.buf, key), f)
	return r
}

func (r *Record) Floats64(key string, f []float64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendFloats64(enc.AppendKey(r.buf, key), f)
	return r
}

func (r *Record) Hex(key string, val []byte) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendHex(enc.AppendKey(r.buf, key), val)
	return r
}

func (r *Record) IPAddr(key string, ip net.IP) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendIPAddr(enc.AppendKey(r.buf, key), ip)
	return r
}

func (r *Record) IPPrefix(key string, pfx net.IPNet) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendIPPrefix(enc.AppendKey(r.buf, key), pfx)
	return r
}

func (r *Record) Int(key string, i int) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInt(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Int16(key string, i int16) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInt16(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Int32(key string, i int32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInt32(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Int64(key string, i int64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInt64(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Int8(key string, i int8) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInt8(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Ints(key string, i []int) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInts(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Ints16(key string, i []int16) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInts16(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Ints32(key string, i []int32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInts32(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Ints64(key string, i []int64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInts64(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Ints8(key string, i []int8) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendInts8(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) MACAddr(key string, ha net.HardwareAddr) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendMACAddr(enc.AppendKey(r.buf, key), ha)
	return r
}

func (r *Record) RawJSON(key string, b []byte) *Record {
	if key == errorKey {
		return r
	}
	r.buf = appendJSON(enc.AppendKey(r.buf, key), b)
	return r
}

func appendJSON(dst []byte, j []byte) []byte {
	return append(dst, j...)
}

func (r *Record) Str(key, val string) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendString(enc.AppendKey(r.buf, key), val)
	return r
}

func (r *Record) Strs(key string, vals []string) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendStrings(enc.AppendKey(r.buf, key), vals)
	return r
}

// Time appends time in [time.RFC3339] format to record.
func (r *Record) Time(key string, t time.Time) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendTime(enc.AppendKey(r.buf, key), t)
	return r
}

func (r *Record) Type(key string, val interface{}) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendType(enc.AppendKey(r.buf, key), val)
	return r
}

func (r *Record) Uint(key string, i uint) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUint(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uint16(key string, i uint16) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUint16(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uint32(key string, i uint32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUint32(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uint64(key string, i uint64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUint64(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uint8(key string, i uint8) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUint8(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uints(key string, i []uint) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUints(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uints16(key string, i []uint16) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUints16(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uints32(key string, i []uint32) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUints32(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uints64(key string, i []uint64) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUints64(enc.AppendKey(r.buf, key), i)
	return r
}

func (r *Record) Uints8(key string, i []uint8) *Record {
	if key == errorKey {
		return r
	}
	r.buf = enc.AppendUints8(enc.AppendKey(r.buf, key), i)
	return r
}

package encoder

import (
	"net"
	"time"

	"github.com/xeptore/flaw/v2/internal/encoder/json"
)

var (
	enc = json.Encoder{}
)

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

func (d *Record) JSON() []byte {
	d.buf = enc.AppendEndMarker(d.buf)
	d.buf = enc.AppendLineBreak(d.buf)
	return d.buf
}

func (d *Record) Err(key string, err error) *Record {
	if err == nil {
		return d
	}
	return d.Str(key, err.Error())
}

func (d *Record) Bool(key string, b bool) *Record {
	d.buf = enc.AppendBool(enc.AppendKey(d.buf, key), b)
	return d
}

func (d *Record) Bools(key string, b []bool) *Record {
	d.buf = enc.AppendBools(enc.AppendKey(d.buf, key), b)
	return d
}

func (d *Record) Dict(key string, dict *Record) *Record {
	dict.buf = enc.AppendEndMarker(dict.buf)
	d.buf = append(enc.AppendKey(d.buf, key), dict.buf...)
	putRecord(dict)
	return d
}

func putRecord(r *Record) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(r.buf) > maxSize {
		return
	}
	//eventPool.Put(e)
}

func (d *Record) Float32(key string, f float32) *Record {
	d.buf = enc.AppendFloat32(enc.AppendKey(d.buf, key), f)
	return d
}

func (d *Record) Float64(key string, f float64) *Record {
	d.buf = enc.AppendFloat64(enc.AppendKey(d.buf, key), f)
	return d
}

func (d *Record) Floats32(key string, f []float32) *Record {
	d.buf = enc.AppendFloats32(enc.AppendKey(d.buf, key), f)
	return d
}

func (d *Record) Floats64(key string, f []float64) *Record {
	d.buf = enc.AppendFloats64(enc.AppendKey(d.buf, key), f)
	return d
}

func (d *Record) Hex(key string, val []byte) *Record {
	d.buf = enc.AppendHex(enc.AppendKey(d.buf, key), val)
	return d
}

func (d *Record) IPAddr(key string, ip net.IP) *Record {
	d.buf = enc.AppendIPAddr(enc.AppendKey(d.buf, key), ip)
	return d
}

func (d *Record) IPPrefix(key string, pfx net.IPNet) *Record {
	d.buf = enc.AppendIPPrefix(enc.AppendKey(d.buf, key), pfx)
	return d
}

func (d *Record) Int(key string, i int) *Record {
	d.buf = enc.AppendInt(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Int16(key string, i int16) *Record {
	d.buf = enc.AppendInt16(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Int32(key string, i int32) *Record {
	d.buf = enc.AppendInt32(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Int64(key string, i int64) *Record {
	d.buf = enc.AppendInt64(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Int8(key string, i int8) *Record {
	d.buf = enc.AppendInt8(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Ints(key string, i []int) *Record {
	d.buf = enc.AppendInts(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Ints16(key string, i []int16) *Record {
	d.buf = enc.AppendInts16(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Ints32(key string, i []int32) *Record {
	d.buf = enc.AppendInts32(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Ints64(key string, i []int64) *Record {
	d.buf = enc.AppendInts64(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Ints8(key string, i []int8) *Record {
	d.buf = enc.AppendInts8(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) MACAddr(key string, ha net.HardwareAddr) *Record {
	d.buf = enc.AppendMACAddr(enc.AppendKey(d.buf, key), ha)
	return d
}

func (d *Record) RawJSON(key string, b []byte) *Record {
	d.buf = appendJSON(enc.AppendKey(d.buf, key), b)
	return d
}

func appendJSON(dst []byte, j []byte) []byte {
	return append(dst, j...)
}

func (d *Record) Str(key, val string) *Record {
	d.buf = enc.AppendString(enc.AppendKey(d.buf, key), val)
	return d
}

func (d *Record) Strs(key string, vals []string) *Record {
	d.buf = enc.AppendStrings(enc.AppendKey(d.buf, key), vals)
	return d
}

func (d *Record) Time(key string, t time.Time) *Record {
	d.buf = enc.AppendTime(enc.AppendKey(d.buf, key), t)
	return d
}

func (d *Record) Type(key string, val interface{}) *Record {
	d.buf = enc.AppendType(enc.AppendKey(d.buf, key), val)
	return d
}

func (d *Record) Uint(key string, i uint) *Record {
	d.buf = enc.AppendUint(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uint16(key string, i uint16) *Record {
	d.buf = enc.AppendUint16(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uint32(key string, i uint32) *Record {
	d.buf = enc.AppendUint32(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uint64(key string, i uint64) *Record {
	d.buf = enc.AppendUint64(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uint8(key string, i uint8) *Record {
	d.buf = enc.AppendUint8(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uints(key string, i []uint) *Record {
	d.buf = enc.AppendUints(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uints16(key string, i []uint16) *Record {
	d.buf = enc.AppendUints16(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uints32(key string, i []uint32) *Record {
	d.buf = enc.AppendUints32(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uints64(key string, i []uint64) *Record {
	d.buf = enc.AppendUints64(enc.AppendKey(d.buf, key), i)
	return d
}

func (d *Record) Uints8(key string, i []uint8) *Record {
	d.buf = enc.AppendUints8(enc.AppendKey(d.buf, key), i)
	return d
}

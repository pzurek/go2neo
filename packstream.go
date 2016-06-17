package go2neo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

type ByteStructures struct {
	typeCode uint8
	size     int
}

// Useful constant definitions
const (
	TinyStringStart = 0x80
	TinyStringEnd   = 0x8F
	StringSize8     = 0xD0
	StringSize16    = 0xD1
	StringSize32    = 0xD2

	TinyListStart = 0x90
	TinyListEnd   = 0x9F
	ListSize8     = 0xD4
	ListSize16    = 0xD5
	ListSize32    = 0xD6

	TinyMapStart = 0xA0
	TinyMapEnd   = 0xAF
	MapSize8     = 0xD8
	MapSize16    = 0xD9
	MapSize32    = 0xDA

	TinyStructStart = 0xB0
	TinyStructEnd   = 0xBF
	StructSize8     = 0xDC
	StructSize16    = 0xDD

	Null    = 0xC0
	Float64 = 0xC1
	False   = 0xC2
	True    = 0xC3
	Int8    = 0xC8
	Int16   = 0xC9
	Int32   = 0xCA
	Int64   = 0xCB

	MinTinyInt = -16
	MaxTinyInt = 127
)

type Decoder struct {
	br *bufio.Reader
}

func NewDecoder(reader io.Reader) (*Decoder, error) {
	dec := Decoder{bufio.NewReader(reader)}
	return &dec, nil
}

type Encoder struct {
	bw *bufio.Writer
}

func NewEncoder(writer io.Writer) (*Encoder, error) {
	enc := Encoder{bufio.NewWriter(writer)}
	return &enc, nil
}

func (dec Decoder) Decode() (interface{}, error) {
	b, err := dec.br.ReadByte()
	if err != nil {
		return nil, err
	}
	if int8(b) >= MinTinyInt && int8(b) <= MaxTinyInt {
		return int(int8(b)), nil
	}
	switch b {
	case False:
		return false, nil
	case True:
		return true, nil
	case Int8:
		i, err := dec.br.ReadByte()
		return int(int8(i)), err
	case Int16:
		barr := [2]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 2 {
			return nil, errors.New("not enough bytes read")
		}
		v := int16(0)
		v |= int16(barr[0]) << 8
		v |= int16(barr[1])
		return int(int16(v)), nil
	case Int32:
		barr := [4]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 4 {
			return nil, errors.New("not enough bytes read")
		}
		v := int32(0)
		v |= int32(barr[0]) << 24
		v |= int32(barr[1]) << 16
		v |= int32(barr[2]) << 8
		v |= int32(barr[3])
		return int(int32(v)), nil
	case Int64:
		barr := [8]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 8 {
			return nil, errors.New("not enough bytes read")
		}
		v := int64(0)
		v |= int64(barr[0]) << 56
		v |= int64(barr[1]) << 48
		v |= int64(barr[2]) << 40
		v |= int64(barr[3]) << 32
		v |= int64(barr[4]) << 24
		v |= int64(barr[5]) << 16
		v |= int64(barr[6]) << 8
		v |= int64(barr[7])
		return int(int64(v)), nil
	}
	return nil, errors.New(fmt.Sprintf("decode: unsupported type: %x", b))
}

func (enc Encoder) Encode(v interface{}) error {
	switch v.(type) {
	case bool:
		return enc.encodeBool(v.(bool))
	case float64:
		return enc.encodeFloat64(v.(float64))
	case int64:
		return enc.encodeInt64(v.(int64))
	case int:
		return enc.encodeInt64(int64(v.(int)))
	}
	return errors.New(fmt.Sprintf("unsupported type: %v", reflect.TypeOf(v)))
}

func (enc Encoder) encodeBool(b bool) error {
	var n int
	var err error
	if b {
		n, err = enc.bw.Write([]byte{True})
	} else {
		n, err = enc.bw.Write([]byte{False})
	}
	if err != nil {
		return err

	}
	if n == 0 {
		return errors.New("failed to encodebool")
	}
	err = enc.bw.Flush()
	return err
}

func (enc Encoder) encodeFloat64(f float64) error {
	v := math.Float64bits(f)
	b := [9]byte{}
	b[0] = Float64
	b[1] = byte(v >> 56)
	b[2] = byte(v >> 48)
	b[3] = byte(v >> 40)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 24)
	b[6] = byte(v >> 16)
	b[7] = byte(v >> 8)
	b[8] = byte(v)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 9 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt64(i int64) error {
	switch {
	case (i >= math.MinInt64 && i < math.MinInt32) || (i > math.MaxInt32 && i <= math.MaxInt64): // INT_64
		b := [9]byte{}
		b[0] = Int64
		b[1] = byte(i >> 56)
		b[2] = byte(i >> 48)
		b[3] = byte(i >> 40)
		b[4] = byte(i >> 32)
		b[5] = byte(i >> 24)
		b[6] = byte(i >> 16)
		b[7] = byte(i >> 8)
		b[8] = byte(i)
		n, err := enc.bw.Write(b[:])
		if err != nil {
			return err
		}
		if n != 9 {
			return errors.New("not enough bytes written")
		}
		enc.bw.Flush()
		return err
	case (i >= math.MinInt32 && i < math.MinInt16) || (i > math.MaxInt16 && i <= math.MaxInt32): // INT_32
		return enc.encodeInt32(int(i))
	case (i >= math.MinInt16 && i < math.MinInt8) || (i > math.MaxInt8 && i <= math.MaxInt16): // INT_16
		return enc.encodeInt16(int(i))
	case i >= math.MinInt8 && i < MinTinyInt: // INT_8
		return enc.encodeInt8(int(i))
	case i >= MinTinyInt && i <= MaxTinyInt: // TINY_INT
		return enc.encodeTinyInt(int(i))
	}
	return errors.New("invalid int64: this should not happen")
}

func (enc Encoder) encodeTinyInt(i int) error {
	if i < MinTinyInt || i > MaxTinyInt {
		return errors.New("encode tinyint: out of range")
	}
	b := byte(i)
	n, err := enc.bw.Write([]byte{b})
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt8(i int) error {
	if i < math.MinInt8 || i > math.MaxInt8 {
		return errors.New("encode int8: out of range")
	}
	n, err := enc.bw.Write([]byte{Int8, byte(i)})
	if err != nil {
		return err
	}
	if n != 2 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt16(i int) error {
	if i < math.MinInt16 || i > math.MaxInt16 {
		return errors.New("encode int16: out of range")
	}
	b := [3]byte{}
	b[0] = Int16
	b[1] = byte(i >> 8)
	b[2] = byte(i)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt32(i int) error {
	if i < math.MinInt32 || i > math.MaxInt32 {
		return errors.New("encode int32: out of range")
	}
	b := [5]byte{}
	b[0] = Int32
	b[1] = byte(i >> 24)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 8)
	b[4] = byte(i)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 5 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

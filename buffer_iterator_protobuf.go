package goutil

import (
	"math"
	"os"
)

const (
	PBWIRE_VARINT           = 0
	PBWIRE_FIXED64          = 1
	PBWIRE_LENGTH_DELIMITED = 2

	PBWIRE_FIXED32 = 5
)

func NewPBKey(field uint, wireType int) *PBKey {
	return &PBKey{
		Field:    field,
		WireType: wireType,
	}
}

type PBKey struct {
	Field    uint
	WireType int
}

func NewPBKeyValue(key PBKey, value []byte) *PBKeyValue {
	return &PBKeyValue{
		Key:   key,
		Value: value,
	}
}

type PBKeyValue struct {
	Key   PBKey
	Value []byte
}

func (iter *BufferIterator) ReadPBKey() (*PBKey, error) {
	if n, err := iter.ReadPBUInt32(); err != nil {
		return nil, err
	} else {
		return NewPBKey(uint(n>>3), int(n&0x07)), nil
	}
}

func (iter *BufferIterator) ReadPBKeyFirstByte(firstByte byte) (*PBKey, error) {
	if firstByte < 128 {
		return NewPBKey(uint(firstByte>>3), int(firstByte&0x07)), nil
	} else {
		n, err := iter.ReadPBUInt32()
		if err != nil {
			return nil, err
		}
		fieldID := uint(n<<4) | uint((firstByte>>3)&0x0F)
		return NewPBKey(fieldID, int(firstByte&0x07)), nil
	}
}
func (iter *BufferIterator) WritePBKey(key *PBKey) error {
	var n uint = (key.Field << 3) | uint(key.WireType)
	return iter.WritePBUInt32(uint32(n))
}
func (iter *BufferIterator) SkipKey(key *PBKey) error {
	switch key.WireType {
	case PBWIRE_FIXED32:
		return iter.Seek(4, os.SEEK_CUR)
	case PBWIRE_FIXED64:
		return iter.Seek(8, os.SEEK_CUR)
	case PBWIRE_LENGTH_DELIMITED:
		n, err := iter.ReadPBUInt32()
		if err != nil {
			return err
		}
		return iter.Seek(int(n), os.SEEK_CUR)
	case PBWIRE_VARINT:
		return iter.ReadPBSkipVarInt()
	default:
		return ErrBufferIteratorNotImplement
	}
}
func (iter *BufferIterator) ReadPBValueBytes(key *PBKey) ([]byte, error) {

	switch key.WireType {
	case PBWIRE_FIXED32:
		return iter.ReadBytes(4)
	case PBWIRE_FIXED64:
		return iter.ReadBytes(8)
	case PBWIRE_LENGTH_DELIMITED:
		nbytes, err := iter.ReadBytes(4)
		if err != nil {
			return nil, err
		}
		subIter := NewBufferIterator(nbytes)
		n, err := subIter.ReadPBUInt32()
		if err != nil {
			return nil, err
		}
		data, err := iter.ReadBytes(int(n))
		if err != nil {
			return nil, err
		}
		ret := []byte{}
		ret = append(ret, nbytes...)
		ret = append(ret, data...)
		return ret, nil
	case PBWIRE_VARINT:
		return iter.ReadPBVarIntBytes()
	default:
		return nil, ErrBufferIteratorNotImplement
	}
}

//MARK: Primary - Protobuff
func (iter *BufferIterator) ReadPBUInt32() (uint32, error) {

	var val uint32 = 0
	for n := 0; n < 5; n++ {

		b, err := iter.ReadByte()

		if err != nil {
			return math.MaxUint32, err
		}

		if (n == 4) && (b&0xF0) != 0 {
			return math.MaxUint32, NewProtocolBufferError("Got larger VarInt than 32bit unsigned")
		}
		if (b & 0x80) == 0 {
			return val | uint32(b<<(7*n)), nil
		}

		val |= (uint32)(b&0x7F) << (7 * n)
	}
	return math.MaxUint32, NewProtocolBufferError("Got larger VarInt than 32bit unsigned")
}

func (iter *BufferIterator) WritePBUInt32(val uint32) error {

	var b byte
	for {
		b = (byte)(val & 0x7F)
		val = val >> 7
		if val == 0 {
			if err := iter.WriteByte(b); err != nil {
				return err
			}
			break
		} else {
			b |= 0x80
			if err := iter.WriteByte(b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (iter *BufferIterator) ReadPBUInt64() (uint64, error) {

	var val uint64 = 0
	for n := 0; n < 10; n++ {

		b, err := iter.ReadByte()

		if err != nil {
			return math.MaxUint32, err
		}

		if (n == 9) && (b&0xFE) != 0 {
			return math.MaxUint32, NewProtocolBufferError("Got larger VarInt than 64bit unsigned")
		}
		if (b & 0x80) == 0 {
			return val | uint64(b<<(7*n)), nil
		}

		val |= (uint64)(b&0x7F) << (7 * n)
	}
	return math.MaxUint32, NewProtocolBufferError("Got larger VarInt than 64bit unsigned")
}
func (iter *BufferIterator) WritePBUInt64(val uint64) error {
	var b byte
	for true {
		b = (byte)(val & 0x7F)
		val = val >> 7
		if val == 0 {
			if err := iter.WriteByte(b); err != nil {
				return err
			}
			break
		} else {
			b |= 0x80
			if err := iter.WriteByte(b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (iter *BufferIterator) ReadPBInt64() (int64, error) {

	if val, err := iter.ReadPBUInt64(); err != nil {

		return math.MaxInt64, err

	} else {

		return int64(val), nil
	}
}

func (iter *BufferIterator) WriteInt64(val int64) error {

	return iter.WritePBUInt64(uint64(val))
}

func (iter *BufferIterator) ReadPBZInt64() (int64, error) {

	if val, err := iter.ReadPBUInt64(); err != nil {

		return math.MaxInt64, err

	} else {

		return (int64)(val>>1) ^ ((int64)(val<<63) >> 63), nil
	}
}

func (iter *BufferIterator) ReadPBInt32() (int32, error) {

	if val, err := iter.ReadPBUInt64(); err != nil {

		return math.MaxInt32, nil

	} else {

		return int32(val), nil
	}
}
func (iter *BufferIterator) WritePBInt32(val int32) error {
	//TODO: check if the integer convert diferrent between C# and golang
	return iter.WritePBUInt64(uint64(val))
}

func (iter *BufferIterator) ReadPBZInt32() (int32, error) {

	if val, err := iter.ReadPBUInt32(); err != nil {

		return math.MaxInt32, err

	} else {

		return (int32)(val>>1) ^ ((int32)(val<<31) >> 31), nil
	}
}

func (iter *BufferIterator) WritePBZInt32(val uint32) error {

	return iter.WritePBUInt32((uint32)((val << 1) ^ (val >> 31)))
}
func (iter *BufferIterator) ReadPBVarIntBytes() ([]byte, error) {
	buff := [10]byte{}
	offset := 0
	for {
		b, err := iter.ReadByte()
		if err != nil {
			return nil, err
		}

		buff[offset] = b
		offset += 1
		if (b & 0x80) == 0 {
			break //end of varint
		}
		if offset >= 10 {
			return nil, NewProtocolBufferError("VarInt too long, more than 10 bytes")
		}
	}
	ret := []byte{}
	ret = append(ret, buff[:offset]...)

	return ret, nil
}
func (iter *BufferIterator) ReadPBSkipVarInt() error {
	for {
		b, err := iter.ReadByte()
		if err != nil {
			return err
		}
		if (b & 0x80) == 0 {
			return nil
		}
	}
}
func (iter *BufferIterator) ReadPBBool() (bool, error) {

	if b, err := iter.ReadByte(); err != nil {
		return true, err
	} else {

		if b == 1 {
			return true, nil
		}
		if b == 0 {
			return false, nil
		}
	}
	return true, NewProtocolBufferError("Invalid boolean value")
}

func (iter *BufferIterator) WritePBBool(val bool) error {

	b := byte(1)
	if !val {
		b = 0
	}
	return iter.WriteByte(b)
}

package goutil

import (
	"encoding/binary"
	"math"
)

func (iter *BufferIterator) ReadFloat32() (float32, error) {
	bytes, err := iter.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float, nil
}
func (iter *BufferIterator) ReadFloat64() (float64, error) {
	bytes, err := iter.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float, nil
}

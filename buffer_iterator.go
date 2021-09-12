package goutil

import (
	"errors"
	"os"
)

var ErrBufferIteratorEOB = errors.New("End of buffer")
var ErrBufferIteratorOutRange = errors.New("Out Of Range")
var ErrBufferIteratorNotImplement = errors.New("Not implemented")

//MARK: Protobuf error
func NewProtocolBufferError(message string) error {
	return errors.New(message)
}

func NewBufferIterator(data []byte) *BufferIterator {

	return &BufferIterator{
		buffer: data,
		offset: 0,
		length: len(data),
	}
}

type BufferIterator struct {
	offset int
	length int
	buffer []byte
}

//EOF return true if reach to the end of buffer
func (iter *BufferIterator) Offset() int {

	return iter.offset
}

func (iter *BufferIterator) Length() int {

	return iter.length
}

func (iter *BufferIterator) EOB() bool {

	return iter.offset >= iter.length
}
func (iter *BufferIterator) Seek(offset int, relative int) error {
	relativeTo := 0
	if relative == os.SEEK_CUR {
		relativeTo = iter.offset
	} else if relative == os.SEEK_END {
		relativeTo = iter.length
	}
	newOffset := relativeTo + offset
	if newOffset < 0 || newOffset >= iter.length {
		return ErrBufferIteratorEOB
	}
	iter.offset = newOffset
	return nil
}

//MARK: Byte
func (iter *BufferIterator) ReadByte() (byte, error) {

	if iter.EOB() {
		return 0xFF, ErrBufferIteratorEOB
	}
	b := iter.buffer[iter.offset]
	iter.offset++
	return b, nil
}
func (iter *BufferIterator) ReadBytes(number int) ([]byte, error) {

	if iter.offset+number >= iter.length {
		return nil, ErrBufferIteratorEOB
	}
	rs := iter.buffer[iter.offset : iter.offset+number]
	iter.offset += number
	return rs, nil
}
func (iter *BufferIterator) WriteByte(b byte) error {
	if iter.EOB() {
		return ErrBufferIteratorEOB
	}
	iter.buffer[iter.offset] = b
	iter.offset++
	return nil
}
func (iter *BufferIterator) WriteBytes(buf []byte) error {

	if iter.offset+len(buf) >= iter.length {

		return ErrBufferIteratorEOB
	}
	for _, b := range buf {

		iter.buffer[iter.offset] = b
		iter.offset++
	}
	return nil
}

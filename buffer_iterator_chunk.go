package goutil

//MARK: Chunk
func (iter *BufferIterator) ReadBBytes() ([]byte, error) {
	return nil, ErrBufferIteratorNotImplement
}
func (iter *BufferIterator) WriteBBytes(buf []byte) error {
	return ErrBufferIteratorNotImplement
}
func (iter *BufferIterator) ReadI32Bytes() ([]byte, error) {
	return nil, ErrBufferIteratorNotImplement
}
func (iter *BufferIterator) WriteI32Bytes(buf []byte) error {
	return ErrBufferIteratorNotImplement
}
func (iter *BufferIterator) ReadI64Bytes() ([]byte, error) {
	return nil, ErrBufferIteratorNotImplement
}
func (iter *BufferIterator) WriteI64Bytes(buf []byte) error {
	return ErrBufferIteratorNotImplement
}

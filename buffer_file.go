package goutil

import (
	"errors"
	"log"
	"os"
	"time"
)

func NewBufferFile(fileName string, data []byte) *BufferFile {

	log.Fatal("not ready yet")
	file := &BufferFile{
		fileName: fileName,
		data:     data,
	}
	return file
}

type BufferFileInfo struct {
	file *BufferFile
}

func (fileInfo *BufferFileInfo) Name() string       { return fileInfo.file.fileName }
func (fileInfo *BufferFileInfo) Size() int64        { return int64(len(fileInfo.file.data)) }
func (fileInfo *BufferFileInfo) Mode() os.FileMode  { return os.ModeIrregular }
func (fileInfo *BufferFileInfo) ModTime() time.Time { return time.Now() }
func (fileInfo *BufferFileInfo) IsDir() bool        { return false }
func (fileInfo *BufferFileInfo) Sys() interface{}   { return nil }

type BufferFile struct {
	fileName string
	data     []byte
	offset   int64
}

//io.Closer
func (f *BufferFile) Close() error {
	return nil
}

//io.Reader
func (f *BufferFile) Read(b []byte) (n int, err error) {
	err = nil
	return
}

//io.Seeker
func (f *BufferFile) Seek(offset int64, whence int) (int64, error) {
	var relativeTo int64
	flen := int64(len(f.data))
	switch whence {
	case 0:
		relativeTo = 0
	case 1:
		relativeTo = f.offset
	case 2:
		relativeTo = flen
	}

	ret := relativeTo + offset
	if ret < 0 || ret > flen {
		return -1, errors.New("Out of file")
	}
	f.offset = ret
	return ret, nil
}

func (f *BufferFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *BufferFile) Stat() (os.FileInfo, error) {
	return nil, nil
}

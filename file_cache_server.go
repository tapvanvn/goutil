package goutil

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

//FILE SERVER
func NewCacheFileServer(fs http.FileSystem) *FileCacheSystem {

	fileSystem := &FileCacheSystem{

		fs: fs,

		cacheFiles: map[string][]byte{},
	}
	return fileSystem
}

type FileCacheSystem struct {
	fs             http.FileSystem
	cacheFiles     map[string][]byte
	totalCacheSize int64
	writeMux       sync.Mutex
}

func (fs *FileCacheSystem) CleanCache() {

	fs.cacheFiles = map[string][]byte{}
}

func (fs *FileCacheSystem) TotalCacheSize() int64 {

	return fs.totalCacheSize
}

func (fs *FileCacheSystem) AddFile(path string, data []byte) {
	fs.writeMux.Lock()
	fs.cacheFiles[path] = data
	fs.writeMux.Unlock()
}

func (fs *FileCacheSystem) RemoveFile(path string) {

	delete(fs.cacheFiles, path)
}

// Open opens file
func (fs FileCacheSystem) Open(path string) (http.File, error) {

	if data, ok := fs.cacheFiles[path]; ok {

		return NewBufferFile(filepath.Base(path), data), nil
	}

	f, err := fs.fs.Open(path)

	if err != nil {

		return nil, err
	}

	s, err := f.Stat()

	if s.IsDir() {

		index := strings.TrimSuffix(path, "/") + "/index.html"

		if _, err := fs.fs.Open(index); err != nil {

			return nil, err
		}

	} else {

		if bytes, err := ioutil.ReadAll(f); err == nil {
			fs.writeMux.Lock()
			fs.cacheFiles[path] = bytes
			fs.writeMux.Unlock()
			f.Close()
			fs.totalCacheSize += int64(len(bytes))
			return NewBufferFile(filepath.Base(path), bytes), nil
		}
	}
	return f, nil
}

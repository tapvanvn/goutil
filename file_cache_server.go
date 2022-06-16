package goutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

type FileServerProxyFunc func(oriPath string) string

//FILE SERVER
func NewCacheFileServer(fs http.FileSystem) *FileCacheSystem {

	fileSystem := &FileCacheSystem{
		Mutex:      &sync.Mutex{},
		fs:         fs,
		prefix:     "/",
		cacheFiles: map[string][]byte{},
		proxies:    map[string]FileServerProxyFunc{},
	}
	return fileSystem
}

type FileCacheSystem struct {
	*sync.Mutex
	fs             http.FileSystem
	cacheFiles     map[string][]byte
	totalCacheSize int64
	prefix         string
	proxies        map[string]FileServerProxyFunc
}

func (fs *FileCacheSystem) ChangeDir(rawfs http.FileSystem) {

	fs.fs = rawfs
	fs.CleanCache()
}

func (fs *FileCacheSystem) CleanCache() {

	fs.cacheFiles = map[string][]byte{}
}

func (fs *FileCacheSystem) TotalCacheSize() int64 {

	return fs.totalCacheSize
}

func (fs *FileCacheSystem) AddFile(path string, data []byte) {
	fs.Lock()
	fs.cacheFiles[path] = data
	fs.Unlock()
}

func (fs *FileCacheSystem) AddProxy(prefix string, fn FileServerProxyFunc) {
	fs.Lock()
	fs.proxies[prefix] = fn
	fs.Unlock()
}

func (fs *FileCacheSystem) RemoveFile(path string) {

	delete(fs.cacheFiles, path)
}

func (fs *FileCacheSystem) SetPrefix(prefix string) {

	fs.prefix = "/" + strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/") + "/"
}

// Open opens file
func (fs FileCacheSystem) Open(path string) (http.File, error) {
	if fs.fs == nil {
		return nil, errors.New("file not found")
	}
	path = "/" + strings.TrimPrefix(path, fs.prefix)

	prefix := path[1:]

	firstSlash := strings.Index(prefix, "/")

	if firstSlash > 0 {

		prefix = prefix[:firstSlash]
	}

	if fn, ok := fs.proxies[prefix]; ok {

		path = fn(strings.TrimPrefix(path, "/"+prefix))
	}
	fs.Lock()

	if data, ok := fs.cacheFiles[path]; ok {
		fs.Unlock()
		return NewBufferFile(filepath.Base(path), data), nil
	}
	fs.Unlock()

	f, err := fs.fs.Open(path)

	if err != nil {
		fmt.Println("err", err.Error())
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
			f.Close()
			fs.Lock()
			fs.cacheFiles[path] = bytes
			fs.Unlock()

			fs.totalCacheSize += int64(len(bytes))
			return NewBufferFile(filepath.Base(path), bytes), nil
		} else {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

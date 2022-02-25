package goutil

import (
	"net/http"
	"strings"
	"sync"
)

func NewFileServer(fs http.FileSystem) *FileSystem {
	fileSystem := &FileSystem{

		Mutex:   &sync.Mutex{},
		FS:      fs,
		prefix:  "/",
		proxies: map[string]FileServerProxyFunc{},
	}
	return fileSystem
}

type FileSystem struct {
	*sync.Mutex
	FS      http.FileSystem
	prefix  string
	proxies map[string]FileServerProxyFunc
}

func (fs *FileSystem) AddProxy(prefix string, fn FileServerProxyFunc) {

	fs.Lock()
	fs.proxies[prefix] = fn
	fs.Unlock()
}

func (fs *FileSystem) SetPrefix(prefix string) {

	fs.prefix = "/" + strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/") + "/"
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {

	path = "/" + strings.TrimPrefix(path, fs.prefix)

	prefix := path[1:]

	firstSlash := strings.Index(prefix, "/")

	if firstSlash > 0 {

		prefix = prefix[:firstSlash]
	}

	if fn, ok := fs.proxies[prefix]; ok {

		path = fn(strings.TrimPrefix(path, "/"+prefix))
	}

	f, err := fs.FS.Open(path)

	if err != nil {

		return nil, err
	}

	s, err := f.Stat()

	if s.IsDir() {

		index := strings.TrimSuffix(path, "/") + "/index.html"

		if _, err := fs.FS.Open(index); err != nil {

			return nil, err
		}
	}

	return f, nil
}

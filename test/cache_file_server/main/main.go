package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tapvanvn/goutil"
)

func proxy(path string) string {
	return "abc" + path
}

func main() {
	rootPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rootPath + "/static")

	cacheFileServer := goutil.NewCacheFileServer(http.Dir(rootPath + "/static"))

	fileServer := http.FileServer(cacheFileServer)

	http.Handle("/", fileServer)

	cacheFileServer2 := goutil.NewCacheFileServer(http.Dir(rootPath + "/static"))

	cacheFileServer2.AddProxy("proxy", proxy)

	cacheFileServer2.SetPrefix("sv2")

	fileServer2 := http.FileServer(cacheFileServer2)

	http.Handle("/sv2/", fileServer2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tapvanvn/goutil"
)

func main() {
	rootPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rootPath + "/static")

	cacheFileServer := goutil.NewCacheFileServer(http.Dir(rootPath + "/static"))

	fileServer := http.FileServer(cacheFileServer)

	http.Handle("/", fileServer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

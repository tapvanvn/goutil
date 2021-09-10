package goutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequestIPAddress(r *http.Request) string {

	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func FromRequest(entity interface{}, r *http.Request) error {

	body, err := ioutil.ReadAll(r.Body)

	fmt.Println(string(body))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &entity)

	if err != nil {
		return err
	}

	return nil
}

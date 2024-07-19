package goutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	content_type := r.Header.Get("Content-Type")

	if content_type == "application/bson" {
		err = bson.Unmarshal(body, entity)
	} else {
		err = json.Unmarshal(body, &entity)
	}

	if err != nil {
		return err
	}

	return nil
}

// ReceiveFile receive the uploaded file
func ReceiveFile(fieldName string, r *http.Request, maxFileLength int64) (string, []byte, error) {

	r.ParseMultipartForm(maxFileLength)

	file, header, err := r.FormFile(fieldName)

	if err != nil {

		return "", nil, err
	}

	defer file.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)

	if err != nil {

		return "", nil, err
	}
	var res = buf.Bytes()

	defer buf.Reset()

	return header.Filename, res, nil
}

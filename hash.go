package goutil

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5 md5 message to footprint
func MD5Message(msg string) string {
	hasher := md5.New()
	hasher.Write([]byte(msg))
	return hex.EncodeToString(hasher.Sum(nil))
}

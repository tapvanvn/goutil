package goutil

import (
	"crypto/md5"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

//MD5 md5 message to footprint
func MD5Message(msg string) string {
	hasher := md5.New()
	return hexutil.Encode(hasher.Sum([]byte(msg)))[2:]
}

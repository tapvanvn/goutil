package goutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

//KeyFromString parse key from string
func KeyFromString(str string) (*rsa.PublicKey, error) {

	pemBlock, _ := pem.Decode([]byte("-----BEGIN RSA PUBLIC KEY-----\n" + str + "\n-----END RSA PUBLIC KEY-----"))

	publicKey, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)

	if err != nil {

		return nil, err

	}

	return publicKey, nil
}

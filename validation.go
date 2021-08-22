package goutil

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"math"
)

//ValidateEmail verify that email is correct format
func ValidateEmail(email string) error {
	return nil
}

//ValidatePassword validate password
func ValidatePassword(password string) error {
	return nil
}

//ValidatePublicKey validate if a public key is valid
func ValidatePublicKey(publicKey string) error {
	return nil
}

//ValidateSignature veriry that signature is correct for content
func ValidateSignature(content []byte, publicKey *rsa.PublicKey, signature string) (bool, error) {

	sig, parseErr := base64.StdEncoding.DecodeString(signature)

	if parseErr != nil {

		return false, parseErr
	}

	h := sha256.New()

	h.Write(content)

	d := h.Sum(nil)

	if verifyErr := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, d, sig); verifyErr == nil {

		return true, nil

	} else {

		return false, verifyErr
	}

}

//IsPowerOfTwo check a number is power of 2
func IsPowerOfTwo(n int) bool {

	if n == 0 {

		return false
	}
	log2 := math.Log2(float64(n))

	return math.Ceil(log2) == math.Floor(log2)
}

package goutil

import (
	"bytes"
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var randArray string = "abcdefghijklmnopqrstuvwxyz0123456789"
var randArrayCase string = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//GenVerifyCode generate a verify code
func GenVerifyCode(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return code
}

//GenOTPCode generate otp code
func GenOTPCode(length int) string {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return strings.ToUpper(code)
}

//GenerateSalt generate a verify code
func GenerateSalt(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArray)
	for i := 0; i < length; i++ {
		code += string(randArray[rand.Intn(arrayLen)])
	}

	return code
}

//GenSecretKey generate secret key
func GenSecretKey(length int) string {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArrayCase)
	for i := 0; i < length; i++ {
		code += string(randArrayCase[rand.Intn(arrayLen)])
	}

	return code
}

//HashPassword Generate hash from password and salt
func HashPassword(password string, salt string) string {
	pwd := []byte(salt + password + salt)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

//VerifyPasswords verify password
func VerifyPasswords(passwordShadow string, password string, salt string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(passwordShadow)
	plainPwd := []byte(salt + password + salt)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func Gen2FASecret() string {
	var b [10]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return base32.StdEncoding.EncodeToString(b[:])
}

//MARK: 2FA - from https://github.com/tilaklodha/google-authenticator
func Prefix2FAOTP(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}
func Get2FAOTP(secret string) (string, error) {

	interval := time.Now().Unix() / 30
	//Converts secret to base32 Encoding. Base32 encoding desires a 32-character
	//subset of the twenty-six letters A–Z and ten digits 0–9
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	//Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)

	if err != nil {
		return "", err
	}
	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000

	//Converts number as a string
	otp := strconv.Itoa(int(h12))

	return Prefix2FAOTP(otp), nil
}

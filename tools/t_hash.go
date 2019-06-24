package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func StrMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func StrSHA256(str string) string {
	hash_value := sha256.New()
	hash_value.Write([]byte(str))
	md := hash_value.Sum(nil)
	return hex.EncodeToString(md)
}

func StrSHA1(str string) string {
	hash_value := sha1.New()
	hash_value.Write([]byte(str))
	md := hash_value.Sum(nil)
	return hex.EncodeToString(md)
}

func StrSHA512(str string) string {
	hash_value := sha512.New()
	hash_value.Write([]byte(str))
	md := hash_value.Sum(nil)
	return hex.EncodeToString(md)
}
